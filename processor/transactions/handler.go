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
	}
	log.Println("applying the processor handler")

	payload := NewMarketPayload(string(request.Payload))
	if payload.IsCreateAccount() {
		handleAccountCreation(payload.CreateAccount(), request.Header, state)
	} else if payload.IsCreateAsset() {
		handleAssetCreation(payload.CreateAsset(), request.Header, state)
	} else if payload.IsCreateHolding() {
		handleHoldingCreation(payload.CreateHolding(), request.Header, state)
	} else if payload.IsCreateOffer() {
		handleOfferCreation(payload.CreateOffer(), request.Header, state)
	} else if payload.IsAcceptOffer() {
		handleOfferAcceptance(payload.AcceptOffer(), request.Header, state)
	} else if payload.IsCloseOffer() {
		handleCloseOffer(payload.CloseOffer(), request.Header, state)
	}
	return nil
}
