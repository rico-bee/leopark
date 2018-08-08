package main

import (
	proto "github.com/golang/protobuf/proto"
	addresser "github.com/rico-bee/leopark/address"
	market "github.com/rico-bee/leopark/market"
	"log"
	"math"
)

type MsgObj interface{}

func containerFromType(space addresser.Space) proto.Message {
	switch space {
	case addresser.SpaceAccount:
		return &market.AccountContainer{}
	case addresser.SpaceAsset:
		return &market.AssetContainer{}
	case addresser.SpaceHolding:
		return &market.HoldingContainer{}
	case addresser.SpaceOffer:
		return &market.OfferContainer{}
	}
	return nil
}

func MapDataToContainer(address string, blockNum int64, data []byte) []MsgObj {
	addressType := addresser.AddressOf(address)
	if addressType == addresser.SpaceOfferHistory {
		return []MsgObj{}
	}
	message := containerFromType(addressType)
	if message == nil {
		log.Println("ignore message from:" + address)
		return nil
	}
	err := proto.Unmarshal(data, message)
	if err != nil {
		log.Println("failed to parse container data:" + err.Error())
	}

	var msgEntries []MsgObj
	switch addressType {
	case addresser.SpaceAccount:
		container := message.(*market.AccountContainer)
		for _, acc := range container.Entries {
			msgEntries = append(msgEntries, mapAccount(blockNum, acc))
		}
	case addresser.SpaceAsset:
		container := message.(*market.AssetContainer)
		for _, ass := range container.Entries {
			msgEntries = append(msgEntries, mapAsset(blockNum, ass))
		}
	case addresser.SpaceHolding:
		container := message.(*market.HoldingContainer)
		for _, h := range container.Entries {
			msgEntries = append(msgEntries, mapHolding(blockNum, h))
		}
	case addresser.SpaceOffer:
		container := message.(*market.OfferContainer)
		for _, o := range container.Entries {
			msgEntries = append(msgEntries, mapOffer(blockNum, o))
		}
	}
	return msgEntries
}

func mapAccount(blockNum int64, account *market.Account) *Account {
	return &Account{
		PublicKey:  account.PublicKey,
		Email:      account.Description,
		Holdings:   account.Holdings,
		BlockRange: BlockRange{StartBlockNum: blockNum, EndBlockNum: math.MaxInt64},
	}
}

func mapAssetRule(rules []*market.Rule) []*Rule {
	mktRules := []*Rule{}
	for _, r := range rules {
		mktRule := &Rule{
			Type:  int32(r.Type),
			Value: string(r.Value),
		}
		mktRules = append(mktRules, mktRule)
	}
	return mktRules
}

func mapAsset(blockNum int64, asset *market.Asset) *Asset {
	return &Asset{
		Name:        asset.Name,
		Description: asset.Description,
		Rules:       mapAssetRule(asset.Rules),
		BlockRange:  BlockRange{StartBlockNum: blockNum, EndBlockNum: math.MaxInt64},
	}
}

func mapHolding(blockNum int64, holding *market.Holding) *Holding {
	return &Holding{
		Id:          holding.Id,
		Label:       holding.Label,
		Description: holding.Description,
		Account:     holding.Account,
		Asset:       holding.Asset,
		Quantity:    holding.Quantity,
		BlockRange:  BlockRange{StartBlockNum: blockNum, EndBlockNum: math.MaxInt64},
	}
}

func mapOffer(blockNum int64, offer *market.Offer) *Offer {
	return &Offer{
		Id:             offer.Id,
		Label:          offer.Label,
		Description:    offer.Description,
		Owners:         offer.Owners,
		Source:         offer.Source,
		SourceQuantity: offer.SourceQuantity,
		Target:         offer.Target,
		TargetQuantity: offer.TargetQuantity,
		Rules:          mapAssetRule(offer.Rules),
		Status:         int32(offer.Status),
		BlockRange:     BlockRange{StartBlockNum: blockNum, EndBlockNum: math.MaxInt64},
	}
}
