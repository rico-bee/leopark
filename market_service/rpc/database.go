package rpc

import (
	"encoding/json"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
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
	_, err = db.DBCreate("market").Run(server.session)
	if err != nil {
		log.Println("failed to create market db" + err.Error())
	}
	_, err = db.DB("market").TableCreate("auth").Run(server.session)
	if err != nil {
		log.Println("failed to create market.auth" + err.Error())
	}
	_, err = db.DB("market").TableCreate("asset").Run(server.session)
	if err != nil {
		log.Println("failed to create market.asset" + err.Error())
	}
	_, err = db.DB("market").TableCreate("holding").Run(server.session)
	if err != nil {
		log.Println("failed to create market.holding" + err.Error())
	}
	_, err = db.DB("market").TableCreate("offer").Run(server.session)
	if err != nil {
		log.Println("failed to create market.offer" + err.Error())
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

func (s *DbServer) CreateAsset(name, description string, rules []*pb.AssetRule) error {
	return db.DB("market").Table("asset").Insert(map[string]interface{}{
		"name":        name,
		"description": description,
		"rules":       rules,
	}).Exec(s.session)
}

func (s *DbServer) CreateHolding(id, label, asset, description string, quantity int64) error {
	return db.DB("market").Table("holding").Insert(map[string]interface{}{
		"label":       label,
		"description": description,
		"asset":       asset,
		"quantity":    quantity,
	}).Exec(s.session)
}

func (s *DbServer) ListAssets() {

}
