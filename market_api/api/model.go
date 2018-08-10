package api

type Block struct {
	BlockId  int64 `gorethink:"block_id"`
	BlockNum int64 `gorethink:"block_num"`
}

type Account struct {
	Email     string     `gorethink:"email" json:"email"`
	PublicKey string     `gorethink:"publicKey" json:"publicKey"`
	Holdings  []*Holding `gorethink:"holdings" json:"holdings"`
}

type Rule struct {
	Type  int32  `gorethink:"type" json:"type"`
	Value string `gorethink:"value" json:"value"`
}

type Asset struct {
	Name        string  `gorethink:"name" json:"name"`
	Description string  `gorethink:"description" json:"description"`
	Rules       []*Rule `gorethink:"rules" json:"rules"`
}

type Holding struct {
	Account  string `gorethink:"account" json:"account"`
	Asset    string `gorethink:"asset" json:"asset"`
	Quantity string `gorethink:"quantity" json:"quantity"`
}
