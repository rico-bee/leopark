package rpc

import (
	"encoding/json"
	db "gopkg.in/gorethink/gorethink.v4"
	"log"
)

type AuthInfo struct {
	Email      string `gorethink:"email"`
	PublicKey  string `gorethink:"publicKey"`
	PwdHash    string `gorethink:"pwdHash,omitempty"`
	PrivateKey string `gorethink:"privateKey,omitempty"`
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
	cursor, err := db.DB("market").Table("auth").Run(s.session)
	if err != nil {
		return nil, err
	}
	var auth AuthInfo
	err = cursor.One(&auth)
	if err != nil {
		log.Println("query error:" + err.Error())
	}
	cursor.Close()
	printObj(auth)
	log.Println(auth.Email)
	return &auth, err
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	log.Println(string(vBytes))
}

func (s *DbServer) CreateUser(authInfo *AuthInfo) error {
	return db.DB("market").Table("auth").Insert(map[string]string{
		"email":      authInfo.Email,
		"pwdHash":    authInfo.PwdHash,
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
