package main

import (
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
		session: session,
	}
	return server, nil
}

func (s *DbServer) fetch(table, id string) (interface{}, error) {
	cursor, err := db.DB(s.Name).Table(table).Run(s.session)
	var row interface{}
	err = cursor.One(&row)
	if err != nil {
		log.Println("query error:" + err.Error())
	}
	cursor.Close()
	return &row, err
}

func (s *DbServer) insert(table string, data interface{}) error {
	return db.DB(s.Name).Table(table).Insert(data).Exec(s.session)
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

func (s *DbServer) DropFork(blockNum int64) (map[string]int64, error) {
	cursor, err := db.DB(s.Name).Table("blocks").Filter(db.Row.Field("block_num").Ge(blockNum)).Delete().Run(s.session)
	if err != nil {
		log.Println("failed to query:" + err.Error())
		return nil, err
	}
	items := map[string]int64{}
	err = cursor.All(&items)

	cursor, err = db.DB(s.Name).TableList().ForEach(func(args ...db.Term) {
		db.Branch(
			db.Eq(args[0], "blocks"),
			[]string{},
			db.Eq(args[0], "auth"),
			[]string{},
			db.DB(s.Name).Table(args[0]).Filter(db.Row.Field("start_block_num").Ge(blockNum)).Delete(),
		)
	}).Run(s.session)

	items2 := map[string]int64{}
	err = cursor.All(&items2)

	ret := map[string]int64{}
	for k, v := range items {
		ret[k] = v + items2[k]
	}
	return ret, nil
}
