package main

import (
	proto "github.com/golang/protobuf/proto"
	events "github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	transaction "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_receipt_pb2"
	addresser "github.com/rico-bee/leopark/address"
	r "gopkg.in/gorethink/gorethink.v4"
	"log"
	"math"
	"regexp"
	"strconv"
)

var (
	NS_REGEX, _ = regexp.Compile("^" + addresser.NS)
)

func processEventList(eventStream <-chan *events.EventList, db *DbServer) {
	eventList := <-eventStream
	blockStream := make(chan *Block)
	blockCommitStream := make(chan *events.Event)
	stateChangeStream := make(chan *events.Event)

	defer close(blockCommitStream)
	defer close(stateChangeStream)
	defer close(blockStream)

	go parseNewBlock(blockCommitStream, blockStream)
	go processStateChange(stateChangeStream, blockStream, db)
	for _, e := range eventList.Events {
		log.Println("processing event:" + e.EventType)
		if e.EventType == "sawtooth/block-commit" {
			blockCommitStream <- e
		} else if e.EventType == "sawtooth/state-delta" {
			stateChangeStream <- e
		}
	}
}

func parseNewBlock(eventStream <-chan *events.Event, blockStream chan<- *Block) {
	block := &Block{}
	e := <-eventStream
	for _, a := range e.Attributes {
		if a.Key == "block_id" {
			block.BlockId = a.Value
		}
		if a.Key == "block_num" {
			blockNum, err := strconv.ParseInt(a.Value, 10, 64)
			if err != nil {
				continue
			}
			block.BlockNum = blockNum
		}
	}
	blockStream <- block
}

func processStateChange(eventStream <-chan *events.Event, blockStream <-chan *Block, db *DbServer) {
	b := <-blockStream
	isDuplicate := resolveIfForked(db, b.BlockNum, b.BlockId)
	if !isDuplicate {
		go stateChanges(eventStream, b, db)
	}
}

func stateChanges(eventStream <-chan *events.Event, block *Block, db *DbServer) {
	e := <-eventStream
	var stateChangeList transaction.StateChangeList
	err := proto.Unmarshal(e.Data, &stateChangeList)
	if err != nil {
		//return nil, errors.New(err.Error())
		log.Println(err.Error())
	}
	states := []*transaction.StateChange{}
	for _, s := range stateChangeList.StateChanges {
		log.Println("checking:" + s.Address)
		if NS_REGEX.MatchString(s.Address) {
			states = append(states, s)
		}
	}
	applyStateChanges(db, states, block.BlockNum)
	insertBlock(db, block.BlockNum, block.BlockId)
}

func resolveIfForked(db *DbServer, blockNum int64, blockId string) bool {
	oldBlk, err := db.fetch("blocks", blockId)
	if err != nil {
		log.Println(err.Error())
	}
	if oldBlk != nil {
		oBlock := oldBlk.(*Block)
		if oBlock.BlockNum == blockNum {
			return true
		}
		ret, err := db.DropFork(blockNum)
		if ret["deleted"] == 0 {
			log.Println("Failed to drop forked resources since block:" + blockId + ":" + err.Error())
		}
	}
	return false
}

func insertBlock(db *DbServer, blockNum int64, blockId string) error {
	newBlock := &Block{
		BlockNum: blockNum,
		BlockId:  blockId,
	}
	return db.insert("blocks", newBlock)
}

func mapAddresSpaceToTable(space addresser.Space) string {
	switch space {
	case addresser.SpaceAccount:
		return "account"
	case addresser.SpaceAsset:
		return "asset"
	case addresser.SpaceHolding:
		return "holding"
	case addresser.SpaceOffer:
		return "offer"
	}
	return ""
}

func findIndex(space addresser.Space, r MsgObj) (string, string) {
	switch space {
	case addresser.SpaceAccount:
		return "public_key", r.(Account).PublicKey
	case addresser.SpaceAsset:
		return "name", r.(Asset).Name
	case addresser.SpaceHolding:
		return "id", r.(Holding).Id
	case addresser.SpaceOffer:
		return "id", r.(Offer).Id
	}
	return "", ""
}

func applyStateChanges(db *DbServer, changes []*transaction.StateChange, blockNum int64) {
	for _, c := range changes {
		resources := MapDataToContainer(c.Address, blockNum, c.Value)
		for _, r := range resources {
			update(db, blockNum, c.Address, r)
		}
	}
}

func update(db *DbServer, blockNum int64, address string, resource MsgObj) (*r.Cursor, error) {
	space := addresser.AddressOf(address)
	table := mapAddresSpaceToTable(space)
	if table == "" {
		log.Println("invalid address detected, cannot update the block")
	}

	idx, idxVal := findIndex(space, resource)
	query := db.Table(table)
	updateQuery := query.GetAll(idxVal, map[string]string{"index": idx}).Filter(r.Row.Field("start_block_num").Eq(math.MaxInt64)).Update(map[string]interface{}{
		"end_block_num": blockNum,
	}).Merge(query.Insert(resource).Without("replaced"))

	return db.Exec(updateQuery)
}
