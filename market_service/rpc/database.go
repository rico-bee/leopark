package rpc

import (
	"log"

	db "gopkg.in/gorethink/gorethink.v4"
)

type AuthInfo struct {
	Email      string `json:"email"`
	PublicKey  string `json:"public_key"`
	PwdHash    string `json:"pwd_hash,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type DbServer struct {
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

func (s *DbServer) FindUser(email string) (*AuthInfo, error) {
	res, err := db.DB("market").Table("auth").Get(email).Run(s.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var auth AuthInfo
	err = res.One(&auth)
	return &auth, err
}

func (s *DbServer) CreateUser(authInfo *AuthInfo) error {
	return db.DB("market").Table("auth").Insert(map[string]string{
		"id":         authInfo.Email,
		"password":   authInfo.PwdHash,
		"privateKey": authInfo.PrivateKey,
		"publicKey":  authInfo.PublicKey,
	}).Exec(s.session)
}

func (s *DbServer) ListUsers(authInfo *AuthInfo) ([]AuthInfo, error) {
	rows, err := db.DB("market").Table("auth").Run(s.session)
	if err != nil {
		return nil, err
	}
	var authInfoList []AuthInfo
	err = rows.All(&authInfoList)
	if err != nil {
		return nil, err
	}
	return authInfoList, nil
}