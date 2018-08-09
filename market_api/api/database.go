package api

import (
	// "encoding/json"
	crypto "github.com/rico-bee/leopark/crypto"
	r "gopkg.in/gorethink/gorethink.v4"
	"log"
	"strconv"
)

type DbServer struct {
	session *r.Session
}

func NewDBServer(url string) (*DbServer, error) {
	session, err := r.Connect(r.ConnectOpts{
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

func (db *DbServer) latestBlockNum() int64 {
	cursor, err := r.DB("market").Table("blocks").Max("block_num").Run(db.session)
	if err != nil {
		log.Println("failed to find latest block number")
		return 0
	}
	var b interface{}
	cursor.One(&b)
	//	blk := b.(Block)
	return 3
}

func (db *DbServer) FindAssets() ([]Asset, error) {
	blkNum := db.latestBlockNum()
	log.Println("latest block:" + strconv.FormatInt(blkNum, 10))
	cursor, err := r.DB("market").Table("asset").Filter(r.Row.Field("start_block_num").Le(blkNum).And(r.Row.Field("end_block_num").Ge(blkNum))).Without("start_block_num", "end_block_num", "").Run(db.session)
	assets := []Asset{}
	err = cursor.All(&assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (db *DbServer) FindAsset(name string) (*Asset, error) {
	blkNum := db.latestBlockNum()
	cursor, err := r.DB("market").Table("asset").Filter(r.Row.Field("start_block_num").Le(blkNum).And(r.Row.Field("end_block_num").Ge(blkNum))).Without("start_block_num", "end_block_num", "").Run(db.session)
	asset := Asset{}
	err = cursor.One(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (s *DbServer) FindUser(email string) (*crypto.AuthInfo, error) {
	cursor, err := r.DB("market").Table("auth").Run(s.session)
	if err != nil {
		return nil, err
	}
	var auth crypto.AuthInfo
	err = cursor.One(&auth)
	if err != nil {
		log.Println("query error:" + err.Error())
	}
	cursor.Close()
	log.Println(auth.Email)
	return &auth, err
}

func (s *DbServer) ListUsers() ([]crypto.AuthInfo, error) {
	rows, err := r.DB("market").Table("auth").Run(s.session)
	if err != nil {
		return nil, err
	}
	var authInfoList []crypto.AuthInfo
	err = rows.All(&authInfoList)
	if err != nil {
		return nil, err
	}
	return authInfoList, nil
}

func (s *DbServer) CreateUser(authInfo *crypto.AuthInfo) error {
	return r.DB("market").Table("auth").Insert(map[string]string{
		"email":      authInfo.Email,
		"pwdHash":    authInfo.PwdHash,
		"privateKey": authInfo.PrivateKey,
		"publicKey":  authInfo.PublicKey,
	}).Exec(s.session)
}

func (s *DbServer) FindAccount(email string) *Account {
	cursor, err := r.DB("market").Table("account").Filter(map[string]interface{}{"email": email}).Merge(func(account r.Term) interface{} {
		holdings := r.DB("market").Table("account").GetAll(map[string]interface{}{"account": account.Field("public_key")})
		return map[string]interface{}{"holdings": holdings}
	}).Run(s.session)
	var acc Account
	err = cursor.One(&acc)
	if err != nil {
		log.Println("failed to find account " + err.Error())
		return nil
	}
	return &acc
}