package transactions

import (
	pb "github.com/rico-bee/marketplace/market"
	proto "github.com/golang/protobuf/proto"
	"log"
)

type MarketPayload struct {
	Transaction *pb.TransactionPayload
}

func NewMarketPayload(payloadStr string) *MarketPayload {
	transaction := &pb.TransactionPayload{}
	err := proto.Unmarshal([]byte(payloadStr), transaction)
	if err != nil {
		log.Fatal(err.Error)
	}
	return &MarketPayload{Transaction: transaction}
}

func (p *MarketPayload) CreateAccount() *pb.CreateAccount {
	return p.Transaction.CreateAccount
}

func (p *MarketPayload) IsCreateAccount() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_CREATE_ACCOUNT
}

func (p *MarketPayload) CreateHolding() *pb.CreateHolding {
	return p.Transaction.CreateHolding
}

func (p *MarketPayload) IsCreateHolding() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_CREATE_HOLDING
}

func (p *MarketPayload) CreateAsset() *pb.CreateAsset {
	return p.Transaction.CreateAsset
}

func (p *MarketPayload) IsCreateAsset() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_CREATE_ASSET
}

func (p *MarketPayload) CreateOffer() *pb.CreateOffer {
	return p.Transaction.CreateOffer
}

func (p *MarketPayload) IsCreateOffer() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_CREATE_OFFER
}

func (p *MarketPayload) AcceptOffer() *pb.AcceptOffer {
	return p.Transaction.AcceptOffer
}

func (p *MarketPayload) IsAcceptOffer() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_ACCEPT_OFFER
}

func (p *MarketPayload) CloseOffer() *pb.CloseOffer {
	return p.Transaction.CloseOffer
}

func (p *MarketPayload) IsCloseOffer() bool {
	return p.Transaction.PayloadType == pb.TransactionPayload_CLOSE_OFFER
}
