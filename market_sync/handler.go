package main

import (
	"encoding/json"
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

func blockParser(done <-chan interface{}, es <-chan *events.Event) <-chan *Block {
	blockStream := make(chan *Block)
	go func() {
		defer close(blockStream)
		var block Block
		e := <-es
		if e.EventType == "sawtooth/block-commit" {
			for _, a := range e.Attributes {
				if a.Key == "block_id" {
					block.BlockId = a.Value
				}
				if a.Key == "block_num" {
					blockNum, err := strconv.ParseInt(a.Value, 10, 64)
					if err != nil {
						log.Println("corrupted block number detected from " + e.EventType)
						return
					}
					block.BlockNum = blockNum
				}
			}
			log.Println("parsing block out:" + block.BlockId)
			blockStream <- &block
		}
	}()
	return blockStream
}

func processEventList(done <-chan interface{}, es <-chan *events.Event, db *DbServer) {
	blockCommitStream := make(chan *events.Event)
	blockDeltaStream := make(chan *events.Event)
	go func() {
		defer close(blockCommitStream)
		defer close(blockDeltaStream)
		for e := range es {
			log.Println("processing:" + e.EventType)
			if e.EventType == "sawtooth/state-delta" {
				select {
				case <-done:
					return
				case blockDeltaStream <- e:
				}
			} else if e.EventType == "sawtooth/block-commit" {
				select {
				case <-done:
					return
				case blockCommitStream <- e:
				}
			}
		}
	}()
	bs := blockParser(done, blockCommitStream)
	processStateChange(done, blockDeltaStream, bs, db)
}

func processStateChange(done <-chan interface{}, deltaStream <-chan *events.Event, blockStream <-chan *Block, db *DbServer) {
	go func() {
		for b := range blockStream {
			log.Println(b.BlockId + "is validated as #" + strconv.FormatInt(b.BlockNum, 10))
			select {
			case <-done:
				return
			case e := <-deltaStream:
				if e != nil && e.EventType == "sawtooth/state-delta" {
					log.Println("block:" + strconv.FormatInt(b.BlockNum, 10) + "-" + b.BlockId)
					isDuplicate := resolveIfForked(db, b.BlockNum, b.BlockId)
					if !isDuplicate {
						stateChanges(e, b, db)
					}
				}
			}
		}
	}()
}

func stateChanges(e *events.Event, block *Block, db *DbServer) {
	var stateChangeList transaction.StateChangeList
	err := proto.Unmarshal(e.Data, &stateChangeList)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("state changes handler" + strconv.Itoa(len(stateChangeList.StateChanges)))
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
	log.Println("calling resolveifforked")
	oldBlk, err := db.fetch("blocks", blockId)
	if err != nil {
		log.Println("cannot fetch old blocks" + err.Error())
	} else {
		if oldBlk != nil {
			jb, err := json.Marshal(oldBlk)
			oBlock := &Block{}
			err = json.Unmarshal(jb, oBlock)
			if oBlock.BlockNum == blockNum {
				return true
			}
			ret, err := db.DropFork(blockNum)
			if ret.Deleted == 0 {
				if err != nil {
					log.Println("query failed:" + err.Error())
				}
				log.Println("Failed to drop forked resources since block: #" + strconv.FormatInt(blockNum, 10))
			}
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
		return "public_key", r.(*Account).PublicKey
	case addresser.SpaceAsset:
		return "name", r.(*Asset).Name
	case addresser.SpaceHolding:
		return "id", r.(*Holding).Id
	case addresser.SpaceOffer:
		return "id", r.(*Offer).Id
	}
	return "", ""
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	log.Println(string(vBytes))
}

func applyStateChanges(db *DbServer, changes []*transaction.StateChange, blockNum int64) {
	log.Println("applying changes")
	for _, c := range changes {
		resources := MapDataToContainer(c.Address, blockNum, c.Value)
		for _, r := range resources {
			log.Println("updating resource:")
			printObj(r)
			_, err := update(db, blockNum, c.Address, r)
			if err != nil {
				log.Println("why we cannot update db:" + err.Error())
			}
		}
	}
}

func update(db *DbServer, blockNum int64, address string, resource MsgObj) (*r.Cursor, error) {
	space := addresser.AddressOf(address)
	table := mapAddresSpaceToTable(space)
	if table == "" {
		log.Println("invalid address detected, cannot update the block")
	}
	log.Println("updating in " + table)
	_, idxVal := findIndex(space, resource)
	query := db.Table(table)
	updateQuery := query.GetAll(idxVal).Filter(r.Row.Field("start_block_num").Eq(math.MaxInt64)).Update(map[string]interface{}{
		"end_block_num": blockNum,
	}).Merge(query.Insert(resource).Without("replaced"))

	return db.Exec(updateQuery)
}
