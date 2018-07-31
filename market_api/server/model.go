package server

type Account struct {
	Email      string `gorethink:"email"`
	PublicKey  string `gorethink:"publicKey"`
	PwdHash    string `gorethink:"pwdHash,omitempty"`
	PrivateKey string `gorethink:"privateKey,omitempty"`
}

type Rule struct {
	Type  int32  `json:"type"`
	Value string `json:"type"`
}

type Asset struct {
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Rules       []*Rule `json: "rules"`
}
