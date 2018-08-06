package transactions

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	addresser "github.com/rico-bee/leopark/address"
	"log"
)

type MarketplaceHandler struct{}

func (h *MarketplaceHandler) FamilyName() string {
	return "LEOTEC"
}

func (h *MarketplaceHandler) FamilyVersions() []string {
	return []string{"1.0"}
}

func (h *MarketplaceHandler) Namespaces() []string {
	return []string{addresser.NS}
}

func (h *MarketplaceHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	state := &MarketState{
		Context: context,
		Timeout: 2,
		State:   make(map[string][]byte),
	}
	log.Println("applying the processor handler")

	payload := NewMarketPayload(string(request.Payload))
	if payload.IsCreateAccount() {
		log.Println("handling account")
		handleAccountCreation(payload.CreateAccount(), request.Header, state)
	} else if payload.IsCreateAsset() {
		log.Println("handling asset")
		handleAssetCreation(payload.CreateAsset(), request.Header, state)
	} else if payload.IsCreateHolding() {
		log.Println("handling holding")
		handleHoldingCreation(payload.CreateHolding(), request.Header, state)
	} else if payload.IsCreateOffer() {
		log.Println("handling offer")
		handleOfferCreation(payload.CreateOffer(), request.Header, state)
	} else if payload.IsAcceptOffer() {
		log.Println("handling accept account")
		handleOfferAcceptance(payload.AcceptOffer(), request.Header, state)
	} else if payload.IsCloseOffer() {
		log.Println("handling close account")
		handleCloseOffer(payload.CloseOffer(), request.Header, state)
	}
	return nil
}
