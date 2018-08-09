package api

type Block struct {
	BlockId  int64 `gorethink:"block_id"`
	BlockNum int64 `gorethink:"block_num"`
}

type Account struct {
	Email     string     `gorethink:"email"`
	PublicKey string     `gorethink:"publicKey"`
	Holdings  []*Holding `gorethink:"holdings"`
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

type Holding struct {
	Account  string `gorethink: "account"`
	Asset    string `gorethink: "asset"`
	Quantity string `gorethink: "quantity"`
}
