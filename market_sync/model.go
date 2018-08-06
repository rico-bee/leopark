package main

type Block struct {
	BlockNum int64  `gorethink:"block_num" json:"block_num"`
	BlockId  string `gorethink:"block_id" json:"block_id"`
	Id       string `gorethink:"id" json:"id,omitempty"`
}

type BlockRange struct {
	StartBlockNum int64 `gorethink:"start_block_num"`
	EndBlockNum   int64 `gorethink:"end_block_num"`
}

type Account struct {
	BlockRange
	Holdings  []string `gorethink:"holdings"`
	PublicKey string   `gorethink:"public_key"`
	Email     string   `gorethink:"email"`
}

type Rule struct {
	Type  int32  `json:"type"`
	Value string `json:"type"`
}

type Asset struct {
	BlockRange
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Rules       []*Rule `json: "rules"`
}

type Holding struct {
	BlockRange
	Id          string `gorethink:"id" json: "id"`
	Label       string `gorethink:"label" json: "label"`
	Account     string `gorethink:"account" json: "account"`
	Description string `gorethink:"description" json: "description"`
	Asset       string `gorethink:"asset" json: "asset"`
	Quantity    int64  `gorethink:"quantity" json: "quantity"`
}

type Offer struct {
	BlockRange
	Id             string   `gorethink:"id" json:"id,omitempty"`
	Label          string   `gorethink:"label" json:"label,omitempty"`
	Description    string   `gorethink:"description" json:"description,omitempty"`
	Owners         []string `gorethink:"owners" json:"owners,omitempty"`
	Source         string   `gorethink:"source" json:"source,omitempty"`
	SourceQuantity int64    `gorethink:"sourceQuantity" json:"source_quantity,omitempty"`
	Target         string   `gorethink:"target" json:"target,omitempty"`
	TargetQuantity int64    `gorethink:"target_quantity" json:"target_quantity,omitempty"`
	Rules          []*Rule  `gorethink:"rules" json:"rules,omitempty"`
	Status         int32    `gorethink:"status" json:"status,omitempty"`
}
