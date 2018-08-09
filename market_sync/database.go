package main

import (
	"encoding/json"
	"errors"
	db "gopkg.in/gorethink/gorethink.v4"
	"log"
)

type DbServer struct {
	Name    string
	session *db.Session
}

func NewDBServer(url string) (*DbServer, error) {
	session, err := db.Connect(db.ConnectOpts{
		Address:    url,
		InitialCap: 10,
		MaxOpen:    10,
	})
	if err != nil {
		log.Println("cannot connect to db")
		return nil, err
	}

	server := &DbServer{
		Name:    "market",
		session: session,
	}

	_, err = db.DBCreate("market").Run(server.session)
	if err != nil {
		log.Println("failed to create market db" + err.Error())
	}
	_, err = db.DB("market").TableCreate("blocks", db.TableCreateOpts{PrimaryKey: "block_num"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.blocks" + err.Error())
	}
	_, err = db.DB("market").TableCreate("account", db.TableCreateOpts{PrimaryKey: "email"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.account" + err.Error())
	}
	_, err = db.DB("market").TableCreate("auth", db.TableCreateOpts{PrimaryKey: "email"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.auth" + err.Error())
	}
	_, err = db.DB("market").TableCreate("asset", db.TableCreateOpts{PrimaryKey: "name"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.asset" + err.Error())
	}
	_, err = db.DB("market").TableCreate("holding", db.TableCreateOpts{PrimaryKey: "id"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.holding" + err.Error())
	}
	_, err = db.DB("market").TableCreate("offer", db.TableCreateOpts{PrimaryKey: "id"}).Run(server.session)
	if err != nil {
		log.Println("failed to create market.offer" + err.Error())
	}
	return server, nil
}

func (s *DbServer) fetch(table, id string) (interface{}, error) {
	cursor, err := db.DB(s.Name).Table(table).Run(s.session)
	var row interface{}
	err = cursor.One(&row)
	if err != nil {
		log.Println("query error:" + err.Error())
		return nil, err
	}
	cursor.Close()
	return &row, err
}

func (s *DbServer) insert(table string, data interface{}) error {
	_, err := db.DB(s.Name).Table(table).Insert(data).RunWrite(s.session)
	if err != nil {
		log.Println("failed to insert data" + err.Error())
	}
	return nil
}

func (s *DbServer) Table(tableName string) db.Term {
	return db.DB(s.Name).Table(tableName)
}

func (s *DbServer) Exec(term db.Term) (*db.Cursor, error) {
	return term.Run(s.session)
}

func (s *DbServer) LastKnownBlocks(numOfBlocks int) ([]string, error) {
	cursor, err := db.DB(s.Name).Table("blocks").OrderBy("block_num").Field("block_id").Run(s.session)
	var rows []string
	err = cursor.All(&rows)
	if err != nil {
		log.Println("failed to get block numbers.")
		return nil, errors.New(err.Error())
	}
	numOfRows := len(rows)
	if numOfRows < numOfBlocks {
		return rows, nil
	}
	return rows[numOfRows-numOfBlocks:], nil
}

type CursorResult struct {
	Deleted   int `json:"deleted"`
	Errors    int `json:"errors"`
	Inserted  int `json:"inserted"`
	Replaced  int `json:"replaced"`
	Skipped   int `json:"skipped"`
	Unchanged int `json:"unchanged"`
}

func unmarshalCursor(v interface{}, ret *CursorResult) {
	vBytes, _ := json.Marshal(v)
	err := json.Unmarshal(vBytes, ret)
	if err != nil {
		ret = nil
	}
}

func (s *DbServer) DropFork(blockNum int64) (*CursorResult, error) {
	dbr, err := db.DB(s.Name).Table("blocks").Filter(db.Row.Field("block_num").Ge(blockNum)).Delete().Run(s.session)
	if err != nil {
		log.Println("failed to query:" + err.Error())
		return nil, err
	}
	ret1 := &CursorResult{}
	unmarshalCursor(dbr, ret1)
	dtr, err := db.DB(s.Name).TableList().ForEach(func(table db.Term) interface{} {
		log.Println("table name:" + table.String())
		return db.Branch(
			db.Eq(table, "blocks"),
			[]string{},
			db.Eq(table, "auth"),
			[]string{},
			db.DB(s.Name).Table(table).Filter(db.Row.Field("start_block_num").Ge(blockNum)).Delete(),
		)
	}).Run(s.session)

	ret2 := &CursorResult{}
	unmarshalCursor(dtr, ret2)
	return &CursorResult{
		Deleted:   ret1.Deleted + ret2.Deleted,
		Errors:    ret1.Errors + ret2.Errors,
		Inserted:  ret1.Inserted + ret2.Inserted,
		Replaced:  ret1.Replaced + ret2.Replaced,
		Skipped:   ret1.Skipped + ret2.Skipped,
		Unchanged: ret1.Unchanged + ret2.Unchanged,
	}, nil
}
