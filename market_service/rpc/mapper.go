package rpc

import (
	mktpb "github.com/rico-bee/leopark/market"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	tx "github.com/rico-bee/leopark/market_service/transaction"
)

func MapHolding(rpcHolding *pb.MarketplaceHolding) *tx.MarketplaceHolding {
	return &tx.MarketplaceHolding{
		HoldingId: rpcHolding.HoldingId,
		Quantity:  rpcHolding.Quantity,
		Asset:     rpcHolding.Asset,
	}
}

func MapAssetRule(rules []*pb.AssetRule) ([]*mktpb.Rule, error) {
	mktRules := []*mktpb.Rule{}
	for _, r := range rules {
		mktRule := &mktpb.Rule{
			Type:  mktpb.Rule_RuleType(r.Type),
			Value: []byte(r.Value),
		}
		mktRules = append(mktRules, mktRule)
	}
	return mktRules, nil
}

func MapOfferParticipant(rpcMktParticipant *pb.OfferParticipant) *tx.OfferParticipant {
	return &tx.OfferParticipant{
		SrcHolding:    rpcMktParticipant.SrcHolding,
		TargetHolding: rpcMktParticipant.TargetHolding,
		SrcAsset:      rpcMktParticipant.SrcAsset,
		TargetAsset:   rpcMktParticipant.TargetAsset,
	}
}
