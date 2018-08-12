package transactions

import (
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	pb "github.com/rico-bee/leopark/market"
	"log"
	"strings"
)

type offerCalcuator struct {
	offer *pb.Offer
	count int64
}

func (c *offerCalcuator) inputQuantity() int64 {
	return c.offer.SourceQuantity * c.count
}

func (c *offerCalcuator) outputQuantity() int64 {
	return c.offer.TargetQuantity * c.count
}

func hasRule(rules []*pb.Rule, ruleType pb.Rule_RuleType) bool {
	for _, rule := range rules {
		if rule.Type == ruleType {
			return true
		}
	}
	return false
}

func hasAccount(account string, accounts []string) bool {
	for _, acc := range accounts {
		if acc == account {
			return true
		}
	}
	return false
}

func ownsAsset(publicKey string, owners []string) bool {
	for _, owner := range owners {
		if owner == publicKey {
			return true
		}
	}
	return false
}

func isNotTransferrable(asset *pb.Asset, ownerKey string) bool {
	if asset == nil {
		log.Println("asset is nil is not transferrable")
	}
	return hasRule(asset.Rules, pb.Rule_NOT_TRANSFERABLE) && !ownsAsset(ownerKey, asset.Owners)
}

func handleOfferCreation(createOffer *pb.CreateOffer, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	offer, err := state.GetOffer(createOffer.Id)
	if err != nil {
		//return nil, errors.New("cannot find the account")
		log.Println("error when getting offer")
	}
	if offer != nil {
		msg := fmt.Sprintf("Offer with id %s already exists", createOffer.Id)
		return nil, errors.New(msg)
	}
	acc, err := state.GetAccount(header.SignerPublicKey)
	if acc == nil {
		//return nil, errors.New("cannot find the account")
		return nil, errors.New("error when getting account")
	}

	log.Println("creating offer......")
	if createOffer.Source == "" {
		return []string{}, errors.New("Cannot find source in createOffer request")
	}
	if createOffer.SourceQuantity == 0 {
		return []string{}, errors.New("Source quantity cannot be 0 in createOffer request")
	}

	log.Println("finding offer source:" + createOffer.Source)
	holding, err := state.GetHolding(createOffer.Source)
	if holding == nil {
		return []string{}, errors.New("Source holding cannot be found")
	}
	if holding.Account != header.SignerPublicKey {
		return []string{}, errors.New("Failed to create Offer, source asset are not transferable")
	}

	if (createOffer.Target != "" && createOffer.TargetQuantity == 0) ||
		(createOffer.Target == "" && createOffer.TargetQuantity > 0) {
		return []string{}, errors.New("failed to create Offer, target and target_quantity must both be set or both unset")
	}

	log.Println("checking target")
	if createOffer.Target != "" {
		targetHolding, err := state.GetHolding(createOffer.Target)
		if err != nil {
			log.Println("cannot find holding")
		}
		if targetHolding == nil {
			return []string{}, errors.New("Failed to create Offer, Holding id listed as target")
		}
		if targetHolding.Account != header.SignerPublicKey {
			return []string{}, errors.New("Failed to create Offer, target Holding account not owned by txn signer")
		}
		targetAsset := state.GetAsset(targetHolding.Asset)
		if targetAsset != nil && isNotTransferrable(targetAsset, header.SignerPublicKey) {
			return []string{}, errors.New("Failed to create Offer, not transferrable")
		}
	}

	log.Println("creating offer with id:" + createOffer.Id)
	return state.SetOffer(
		createOffer.Id,
		createOffer.Label,
		createOffer.Description,
		createOffer.Source,
		createOffer.Target,
		[]string{header.SignerPublicKey},
		createOffer.SourceQuantity,
		createOffer.TargetQuantity,
		createOffer.Rules,
	)
}

func handleCloseOffer(closeOffer *pb.CloseOffer, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	offer, err := state.GetOffer(closeOffer.Id)
	if err != nil {
		log.Println("cannot find the offer")
	}

	if offer.Status != pb.Offer_OPEN {
		return []string{}, errors.New("Offer is not open")
	}

	if !ownsAsset(header.SignerPublicKey, offer.Owners) {
		return []string{}, errors.New("User doesn't own the offer")
	}

	return state.CloseOffer(offer.Id)
}

type OfferParticipant struct {
	SrcHolding    *pb.Holding
	TargetHolding *pb.Holding
	SrcAsset      *pb.Asset
	TargetAsset   *pb.Asset
}

type OfferAcceptance struct {
	offer       *pb.Offer
	acceptOffer *pb.AcceptOffer
	state       *MarketState
	header      *pb2.TransactionHeader
	offerer     *OfferParticipant
	receiver    *OfferParticipant
}

func newOfferAcceptance(acceptOffer *pb.AcceptOffer, header *pb2.TransactionHeader, state *MarketState) *OfferAcceptance {
	offer, err := state.GetOffer(acceptOffer.Id)
	if err != nil {
		return nil
	}

	srcHolding, err := state.GetHolding(offer.Source)
	if err != nil {
		log.Println(err.Error)
	}
	var targetHolding *pb.Holding
	if offer.Target != "" {
		targetHolding, err = state.GetHolding(offer.Target)
		if err != nil {
			log.Println(err.Error)
		}
	}
	var srcAsset *pb.Asset
	if targetHolding != nil {
		srcAsset = state.GetAsset(srcHolding.Asset)
		if err != nil {
			log.Println(err.Error)
		}
	}

	var targetAsset *pb.Asset
	if targetHolding != nil {
		targetAsset = state.GetAsset(targetHolding.Asset)
	}

	offerer := &OfferParticipant{
		SrcHolding:    srcHolding,
		TargetHolding: targetHolding,
		SrcAsset:      srcAsset,
		TargetAsset:   targetAsset,
	}

	var source *pb.Holding
	source, err = state.GetHolding(acceptOffer.Source)
	if err != nil {
		log.Println("")
	}

	var target *pb.Holding
	target, err = state.GetHolding(acceptOffer.Target)
	if err != nil {
		log.Println("")
	}

	var acceptSrcAsset *pb.Asset
	if source != nil {
		acceptSrcAsset = state.GetAsset(source.Asset)
	}

	var acceptTargetAsset *pb.Asset
	if target != nil {
		acceptTargetAsset = state.GetAsset(target.Asset)
	}

	receiver := &OfferParticipant{
		SrcHolding:    source,
		TargetHolding: target,
		SrcAsset:      acceptSrcAsset,
		TargetAsset:   acceptTargetAsset,
	}

	return &OfferAcceptance{
		offer:       offer,
		acceptOffer: acceptOffer,
		state:       state,
		header:      header,
		offerer:     offerer,
		receiver:    receiver,
	}
}

func (a *OfferAcceptance) validateOutputHoldingExists() error {
	if a.offer.Target != "" && a.acceptOffer.Source != "" {
		if a.offerer.TargetHolding != nil && a.receiver.SrcHolding == nil {
			return errors.New("whatever it is")
		}
	}
	return nil
}

func (a *OfferAcceptance) validateInputHoldingExists() error {
	if a.receiver.TargetHolding == nil {
		return errors.New("Invalid target")
	}
	return nil
}

func (a *OfferAcceptance) validateInputHoldingAssets() error {
	if a.offerer.SrcHolding.Asset != a.receiver.TargetHolding.Asset {
		return errors.New("Failed to accept offer, expected Holding asset")
	}
	return nil
}

func (a *OfferAcceptance) validateOutputHoldingAssets() error {
	if a.offer.Target != "" && a.offerer.TargetHolding != nil && a.offerer.TargetHolding.Asset != a.receiver.SrcHolding.Asset {
		return errors.New("Failed to accept offer, expected Holding asset")
	}
	return nil
}

func (a *OfferAcceptance) validateOutputEnough(outputQuantity int64) error {
	if a.acceptOffer.Source != "" && outputQuantity > a.receiver.SrcHolding.Quantity {
		return errors.New("output is not enough")
	}
	return nil
}

func (a *OfferAcceptance) validateInputEnough(inputQuantity int64) error {
	if a.acceptOffer.Source != "" && inputQuantity > a.offerer.SrcHolding.Quantity {
		return errors.New("input is not enough")
	}
	return nil
}

func accountsLimitedTo(offer *pb.Offer) bool {
	return hasRule(offer.Rules, pb.Rule_EXCHANGE_LIMITED_TO_ACCOUNTS)
}

func exchangeOnce(offer *pb.Offer) bool {
	return hasRule(offer.Rules, pb.Rule_EXCHANGE_ONCE)
}

func exchangeOncePerAccount(offer *pb.Offer) bool {
	return hasRule(offer.Rules, pb.Rule_EXCHANGE_ONCE_PER_ACCOUNT)
}

func isHoldingInfinite(asset *pb.Asset, owner string) bool {
	return asset != nil && (hasRule(asset.Rules, pb.Rule_ALL_HOLDINGS_INFINITE) ||
		hasRule(asset.Rules, pb.Rule_OWNER_HOLDINGS_INFINITE))
}

func accounts(offer *pb.Offer) []string {
	accountMap := make(map[string]int)
	for _, r := range offer.Rules {
		if r.Type == pb.Rule_EXCHANGE_LIMITED_TO_ACCOUNTS {
			tokens := strings.Split(string(r.Value), ",")
			if len(tokens) > 0 {
				for _, acc := range tokens {
					accountMap[acc]++
				}
			}
		}
	}
	accounts := []string{}
	for acc, count := range accountMap {
		if count > 1 {
			accounts = append(accounts, acc)
		}
	}
	return accounts
}

func (a *OfferAcceptance) validateOncePerAccount() error {
	if exchangeOncePerAccount(a.offer) {
		if a.state.OfferHasReceipt(a.offer.Id) {
			return errors.New("Failed to accept offer, EXCHANGE ONCE PER ACCOUNT set and account already has accepted offer")
		}
	}
	return nil
}

func (a *OfferAcceptance) validateExchangeOnce() error {
	if exchangeOnce(a.offer) {
		if a.state.OfferHasReceipt(a.offer.Id) {
			return errors.New("offer has been used")
		}
	}
	return nil
}

func (a *OfferAcceptance) validateAccountsLimitedTo() error {
	if accountsLimitedTo(a.offer) && hasAccount(a.header.SignerPublicKey, accounts(a.offer)) {
		return errors.New("offer cannot proceeded")
	}
	return nil
}

func (a *OfferAcceptance) handleOffererSource(inputQuantity int64) ([]string, error) {
	if isHoldingInfinite(a.offerer.SrcAsset, a.offerer.SrcHolding.Account) {
		return a.state.UpdateHolding(a.offerer.SrcHolding.Id, inputQuantity)
	}
	return []string{}, errors.New("offer cannot be processed")
}

func (a *OfferAcceptance) handleOffererTarget(inputQuantity int64) ([]string, error) {
	if a.offer.Target != "" {
		return a.state.UpdateHolding(a.offerer.TargetHolding.Id, inputQuantity)
	}
	return []string{}, errors.New("offer cannot be processed")
}

func (a *OfferAcceptance) handleReceiverSource(outputQuantity int64) ([]string, error) {
	if a.acceptOffer.Source != "" && !isHoldingInfinite(a.receiver.SrcAsset, a.receiver.SrcHolding.Account) {
		return a.state.UpdateHolding(a.receiver.SrcHolding.Id, outputQuantity)
	}
	return []string{}, errors.New("offer cannot be processed")
}

func (a *OfferAcceptance) handleReceiverTarget(inputQuantity int64) ([]string, error) {
	return a.state.UpdateHolding(a.receiver.TargetHolding.Id, inputQuantity)
}

func (a *OfferAcceptance) handleOncePerAccount() ([]string, error) {
	if exchangeOncePerAccount(a.offer) {
		return a.state.SaveOfferAccountReceipt(a.offer.Id, a.header.SignerPublicKey)
	}
	return []string{}, errors.New("holding cannot be processed")
}

func (a *OfferAcceptance) handleExchangeOnce() ([]string, error) {
	if exchangeOnce(a.offer) {
		return a.state.SaveOfferReceipt(a.offer.Id)
	}
	return []string{}, errors.New("holding cannot be processed")
}

func handleOfferAcceptance(acceptOffer *pb.AcceptOffer, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	offer, err := state.GetOffer(acceptOffer.Id)
	if err != nil {
		log.Println("")
		return []string{}, errors.New("User doesn't own the offer")
	}
	if offer == nil || offer.Status != pb.Offer_OPEN {
		return []string{}, errors.New("Offer is not valid")
	}

	offerAcceptance := newOfferAcceptance(acceptOffer, header, state)
	err = offerAcceptance.validateExchangeOnce()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateOncePerAccount()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateAccountsLimitedTo()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateOutputHoldingExists()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateInputHoldingExists()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateInputHoldingAssets()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateOutputHoldingAssets()
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}

	calculator := &offerCalcuator{
		offer: offer,
		count: int64(acceptOffer.Count),
	}

	err = offerAcceptance.validateOutputEnough(calculator.outputQuantity())
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}
	err = offerAcceptance.validateInputEnough(calculator.inputQuantity())
	if err != nil {
		return []string{}, errors.New("offer cannot be processed")
	}

	ret, err := offerAcceptance.handleOffererSource(calculator.inputQuantity())
	if err != nil {
		return ret, err
	}
	ret, err = offerAcceptance.handleOffererTarget(calculator.outputQuantity())
	if err != nil {
		return ret, err
	}
	ret, err = offerAcceptance.handleReceiverSource(calculator.outputQuantity())
	if err != nil {
		return ret, err
	}
	ret, err = offerAcceptance.handleReceiverTarget(calculator.inputQuantity())
	if err != nil {
		return ret, err
	}
	ret, err = offerAcceptance.handleOncePerAccount()
	if err != nil {
		return ret, err
	}
	return offerAcceptance.handleExchangeOnce()
}
