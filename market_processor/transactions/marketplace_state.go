package transactions

import (
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	addresser "github.com/rico-bee/leopark/address"
	pb "github.com/rico-bee/leopark/market"
	"log"
)

type StateEntries map[string][]byte

type MarketState struct {
	Context *processor.Context
	Timeout int
	State   StateEntries
}

func (s *MarketState) FindState(address string) ([]byte, error) {
	return []byte{}, nil
}

func (s *MarketState) GetAccountContainer(address string) (*pb.AccountContainer, error) {
	container := &pb.AccountContainer{}
	if val, ok := s.State[address]; ok {
		err := proto.Unmarshal(val, container)
		if err != nil {
			log.Println("corrupted account container data:" + err.Error())
		}
		return container, nil
	}
	return container, errors.New("Failed to get account container")
}

func (s *MarketState) FindAccountFromContainer(identifier string, container *pb.AccountContainer) (*pb.Account, error) {
	if container == nil || container.Entries == nil {
		return nil, errors.New("account container is nil:" + identifier)
	}
	for _, account := range container.Entries {
		if account.PublicKey == identifier {
			return account, nil
		}
	}
	return nil, errors.New("Cannot find account")
}

func (s *MarketState) GetAssetContainer(address string) (*pb.AssetContainer, error) {
	if val, ok := s.State[address]; ok {
		container := &pb.AssetContainer{}
		proto.Unmarshal(val, container)
		return container, nil
	}
	return nil, errors.New("Failed to get asset container")
}

func (s *MarketState) FindAssetFromContainer(name string, container *pb.AssetContainer) (*pb.Asset, error) {
	if container == nil || container.Entries == nil {
		return nil, errors.New("Cannot find asset")
	}
	for _, asset := range container.Entries {
		if asset.Name == name {
			return asset, nil
		}
	}
	return nil, errors.New("Cannot find asset")
}

func (s *MarketState) GetHoldingContainer(address string) (*pb.HoldingContainer, error) {
	if val, ok := s.State[address]; ok {
		container := &pb.HoldingContainer{}
		proto.Unmarshal(val, container)
		return container, nil
	}
	return nil, errors.New("Failed to get holding container")
}

func (s *MarketState) FindHoldingFromContainer(id string, container *pb.HoldingContainer) (*pb.Holding, error) {
	if container == nil || container.Entries == nil {
		return nil, errors.New("Cannot find holding")
	}
	for _, holding := range container.Entries {
		if holding.Id == id {
			return holding, nil
		}
	}
	return nil, errors.New("Cannot find holding")
}

func (s *MarketState) GetOfferContainer(address string) (*pb.OfferContainer, error) {
	if val, ok := s.State[address]; ok {
		container := &pb.OfferContainer{}
		proto.Unmarshal(val, container)
		return container, nil
	}
	return nil, errors.New("Failed to get offer container")
}

func (s *MarketState) FindOfferFromContainer(id string, container *pb.OfferContainer) (*pb.Offer, error) {
	if container == nil || container.Entries == nil {
		return nil, errors.New("Cannot find offer")
	}
	for _, offer := range container.Entries {
		if offer.Id == id {
			return offer, nil
		}
	}
	return nil, errors.New("Cannot find offer")
}

func (s *MarketState) GetHistoryContainer(address string) (*pb.OfferHistoryContainer, error) {
	if val, ok := s.State[address]; ok {
		container := &pb.OfferHistoryContainer{}
		proto.Unmarshal(val, container)
		return container, nil
	}
	return nil, errors.New("Failed to get offer history container")
}

func (s *MarketState) FindHistoryByOfferId(offerId string, container *pb.OfferHistoryContainer) (*pb.OfferHistory, error) {
	if container == nil || container.Entries == nil {
		return nil, errors.New("Cannot find offer")
	}
	for _, history := range container.Entries {
		if history.OfferId == offerId {
			return history, nil
		}
	}
	return nil, errors.New("Failed to find offer from history")
}

func (s *MarketState) FindHistory(offerId string, accountId string, container *pb.OfferHistoryContainer) (*pb.OfferHistory, error) {
	for _, history := range container.Entries {
		if history.OfferId == offerId && history.AccountId == accountId {
			return history, nil
		}
	}
	return nil, errors.New("Failed to find offer from history")
}

func (s *MarketState) FindOfferAccountReceipt(offerId, account string) (*pb.OfferHistory, error) {
	address := addresser.MakeOfferAccountAddress(offerId, account)
	state, err := s.Context.GetState([]string{address})
	if err != nil {
		log.Println("Failed to get state")
	}
	for address, data := range state {
		s.State[address] = data
	}
	historyContainer, err := s.GetHistoryContainer(address)
	if err != nil {
		log.Println("Failed to get offer history container")
		return nil, errors.New("Failed to find offer history container: " + address)
	}
	offerHistory, err := s.FindHistory(offerId, account, historyContainer)
	if err != nil {
		log.Println("Failed to find offer history")
		return nil, errors.New("Failed to find offer history: " + address)
	}
	return offerHistory, nil
}

func (s *MarketState) OfferHasReceipt(offerId string) bool {
	address := addresser.MakeOfferHistoryAddress(offerId)
	state, err := s.Context.GetState([]string{address})
	if err != nil {
		log.Println("Failed to get state")
		return false
	}
	for address, data := range state {
		s.State[address] = data
	}
	historyContainer, err := s.GetHistoryContainer(address)
	if err != nil {
		log.Println("Failed to find history container")
		return false
	}
	_, err = s.FindHistoryByOfferId(offerId, historyContainer)
	if err != nil {
		return false
	}
	return true
}

func (s *MarketState) SaveOfferReceipt(offerId string) ([]string, error) {
	address := addresser.MakeOfferHistoryAddress(offerId)

	offerHistory := &pb.OfferHistory{
		OfferId: offerId,
	}
	container, err := s.GetHistoryContainer(address)
	if container == nil {
		container = &pb.OfferHistoryContainer{}
	}
	if container.GetEntries() == nil {
		container.Entries = []*pb.OfferHistory{}
	}
	container.Entries = append(container.GetEntries(), offerHistory)
	newState := make(map[string][]byte)
	newState[address], err = proto.Marshal(container)
	if err != nil {
		log.Println("Failed to serialize container")
		return nil, err
	}
	return s.Context.SetState(newState)
}

func (s *MarketState) SaveOfferAccountReceipt(offerId, accountId string) ([]string, error) {
	address := addresser.MakeOfferAccountAddress(offerId, accountId)
	container, err := s.GetHistoryContainer(address)

	if container == nil {
		container = &pb.OfferHistoryContainer{}
	}
	if container.GetEntries() == nil {
		container.Entries = []*pb.OfferHistory{}
	}
	offerHistory := &pb.OfferHistory{
		OfferId:   offerId,
		AccountId: accountId,
	}
	container.Entries = append(container.GetEntries(), offerHistory)
	newState := make(map[string][]byte)
	newState[address], err = proto.Marshal(container)
	if err != nil {
		log.Println("Failed to serialize container")
		return nil, err
	}
	return s.Context.SetState(newState)
}

func (s *MarketState) AddHoldingToAccount(accountId, holdingId string) ([]string, error) {
	address := addresser.MakeAccountAddress(accountId)
	container, err := s.GetAccountContainer(address)
	if container == nil {
		log.Println("cannot find account container:" + err.Error())
		return nil, err
	}

	account, err := s.FindAccountFromContainer(accountId, container)
	// if account == nil {
	// 	account = &pb.Account{
	// 		PublicKey:   accountId,
	// 		Label:       "",
	// 		Description: "",
	// 		Holdings:    []string{},
	// 	}
	// 	if container.GetEntries() == nil {
	// 		container.Entries = []*pb.Account{}
	// 	}
	// 	container.Entries = append(container.GetEntries(), account)
	// }
	if account == nil {
		log.Println("add holding to account failed: cannot find account:" + accountId)
		return nil, errors.New("add holding to account failed: cannot find account:" + accountId)
	}
	account.Holdings = append(account.Holdings, holdingId)
	newState := make(map[string][]byte)
	bytes, err := proto.Marshal(container)
	if err != nil {
		log.Fatal("corrupted container data")
	}
	newState[address] = bytes
	log.Println("writing holdings into account back to chain")
	return s.Context.SetState(newState)
}

func (s *MarketState) SetAccount(accountKey, label, description string, holdings []string) ([]string, error) {

	log.Println("processor: create account: for " + label + " with key " + accountKey)
	address := addresser.MakeAccountAddress(accountKey)
	container, err := s.GetAccountContainer(address)
	account, err := s.FindAccountFromContainer(accountKey, container)
	if account == nil {
		account = &pb.Account{
			PublicKey:   accountKey,
			Label:       label,
			Description: description,
			Holdings:    holdings,
		}
		container.Entries = append(container.Entries, account)
	}
	newState := make(map[string][]byte)
	bytes, err := proto.Marshal(container)
	if err != nil {
		log.Println("corrupted container data")
	}
	newState[address] = bytes
	return s.Context.SetState(newState)
}

func (s *MarketState) fetchState(address string) error {
	if _, ok := s.State[address]; !ok {
		state, err := s.Context.GetState([]string{address})
		if err != nil {
			log.Println("failed to find state for:" + address + ": " + err.Error())
			return err
		}
		for addr, data := range state {
			s.State[addr] = data
		}
	}
	return nil
}

func (s *MarketState) GetAccount(accountKey string) (*pb.Account, error) {
	address := addresser.MakeAccountAddress(accountKey)
	err := s.fetchState(address)
	if err != nil {
		log.Println("cannot fetch state data for account:")
		return nil, err
	}

	container, err := s.GetAccountContainer(address)
	if err != nil {
		log.Println("cannot find account container, empty container returned: " + err.Error())
	}
	return s.FindAccountFromContainer(accountKey, container)
}

func (s *MarketState) SetAsset(name string, description string, owners []string, rules []*pb.Rule) ([]string, error) {
	address := addresser.MakeAssetAddress(name)
	container, err := s.GetAssetContainer(address)
	if container == nil {
		log.Println("cannot get asset container")
		if err != nil {
			log.Println("error: " + err.Error())
		}
		container = &pb.AssetContainer{}
	}

	asset := &pb.Asset{
		Name:        name,
		Description: description,
		Owners:      owners,
		Rules:       rules,
	}
	if container.GetEntries() == nil {
		container.Entries = []*pb.Asset{}
	}
	container.Entries = append(container.Entries, asset)
	newState := make(map[string][]byte)
	bytes, err := proto.Marshal(container)
	newState[address] = bytes
	log.Println("processor: create asset: " + name + " - " + owners[0])
	return s.Context.SetState(newState)
}

func (s *MarketState) GetAsset(name string) *pb.Asset {
	address := addresser.MakeAssetAddress(name)
	err := s.fetchState(address)
	if err != nil {
		log.Println("cannot fetch state data for asset: " + name + ":" + address)
		return nil
	}
	container, err := s.GetAssetContainer(address)
	if err != nil {
		log.Println("cannot get container for : " + name + "due to:" + err.Error())
		return nil
	}
	asset, err := s.FindAssetFromContainer(name, container)
	if err != nil {
		log.Println("cannot find asset " + name)
		return nil
	}
	return asset
}

func (s *MarketState) UpdateHolding(holdingId string, quantity int64) ([]string, error) {
	address := addresser.MakeHoldingAddress(holdingId)
	container, err := s.GetHoldingContainer(address)
	if err != nil {
		log.Println("cannot find holding container for :" + holdingId)
		return nil, err
	}
	holding, err := s.FindHoldingFromContainer(holdingId, container)
	if err != nil {
		log.Println("cannot find the holding")
		return nil, err
	}
	holding.Quantity = quantity
	state := make(map[string][]byte)
	state[address], err = proto.Marshal(container)
	if err != nil {
		log.Println("cannot serialize state container:" + err.Error())
		return nil, err
	}
	return s.Context.SetState(state)
}

func (s *MarketState) CreateHolding(identifier string, label string,
	description string, account string, asset string, quantity int64) ([]string, error) {
	address := addresser.MakeHoldingAddress(identifier)
	container, err := s.GetHoldingContainer(address)
	if err != nil {
		log.Println("cannot get holding container")
	}
	if container == nil {
		container = &pb.HoldingContainer{}
	}
	if container.GetEntries() == nil {
		container.Entries = []*pb.Holding{}
	}
	holding := &pb.Holding{
		Id:          identifier,
		Label:       label,
		Description: description,
		Account:     account,
		Asset:       asset,
		Quantity:    quantity,
	}
	container.Entries = append(container.GetEntries(), holding)
	state := make(map[string][]byte)
	state[address], err = proto.Marshal(container)
	if err != nil {
		log.Println("cannot serialize state container:" + err.Error())
		return nil, err
	}
	return s.Context.SetState(state)
}

func (s *MarketState) GetHolding(identifier string) (*pb.Holding, error) {
	address := addresser.MakeHoldingAddress(identifier)
	err := s.fetchState(address)
	if err != nil {
		log.Println("cannot fetch state data for holding:" + address)
		return nil, err
	}
	container, err := s.GetHoldingContainer(address)
	if err != nil {
		log.Println("failed to find container")
		return nil, err
	}
	return s.FindHoldingFromContainer(identifier, container)
}

func (s *MarketState) CloseOffer(offerId string) ([]string, error) {
	address := addresser.MakeOfferAddress(offerId)
	container, err := s.GetOfferContainer(address)
	if err != nil {
		log.Println("Failed to get container")
		return nil, err
	}
	offer, err := s.FindOfferFromContainer(offerId, container)
	if err != nil {
		log.Println("cannot find the offer")
		return nil, err
	}
	offer.Status = pb.Offer_CLOSED
	state := make(map[string][]byte)
	state[address], err = proto.Marshal(container)
	if err != nil {
		log.Println("cannot serialize state container:" + err.Error())
		return nil, err
	}
	return s.Context.SetState(state)
}

func (s *MarketState) FindOfferRules(holdingId string) []*pb.Rule {
	//addr := addresser.MakeHoldingAddress(holdingId)
	holding, err := s.GetHolding(holdingId)
	if err != nil {
		log.Println("cannot get holding")
		return nil
	}
	asset := s.GetAsset(holding.Asset)

	log.Println("finding offer rule: " + holding.Asset)

	if asset == nil {
		log.Println("failed to find asset")
		return nil
	}
	log.Println("found it:" + asset.Name)
	rules := make([]*pb.Rule, len(asset.Rules), (cap(asset.Rules)+1)*2)
	for idx, r := range asset.Rules {
		rules[idx] = r
	}
	return rules
}

func (s *MarketState) SetOffer(identifier, label, description,
	source, target string, owners []string,
	source_quantity, target_quantity int64,
	rules []*pb.Rule) ([]string, error) {

	addr := addresser.MakeOfferAddress(identifier)
	container, err := s.GetOfferContainer(addr)
	if err != nil {
		log.Println("offer container doesn't exist: " + addr)
		if container == nil {
			log.Println("initializing offer container")
			container = &pb.OfferContainer{}
		}
	}

	offer := &pb.Offer{
		Id:             identifier,
		Label:          label,
		Description:    description,
		Owners:         owners,
		Source:         source,
		SourceQuantity: source_quantity,
		Target:         target,
		TargetQuantity: target_quantity,
		Rules:          rules,
		Status:         pb.Offer_OPEN,
	}

	srcRule := s.FindOfferRules(source)
	if srcRule == nil {
		log.Println("failed to find offer rule for source")
	} else {
		offer.Rules = append(offer.Rules, srcRule...)
	}
	if target != "" {
		targetRule := s.FindOfferRules(target)
		if targetRule == nil {
			log.Println("failed to find offer rule")
		} else {
			offer.Rules = append(offer.Rules, targetRule...)
		}
	}
	log.Println("adding offer into container:")
	container.Entries = append(container.Entries, offer)
	state := make(map[string][]byte)
	state[addr], err = proto.Marshal(container)
	if err != nil {
		log.Println("cannot serialize state container:" + err.Error())
		return nil, err
	}
	return s.Context.SetState(state)
}

func (s *MarketState) GetOffer(identifier string) (*pb.Offer, error) {
	addr := addresser.MakeOfferAddress(identifier)
	err := s.fetchState(addr)
	if err != nil {
		log.Println("cannot fetch state data for offer:" + addr)
		return nil, err
	}
	container, err := s.GetOfferContainer(addr)
	if err != nil {
		log.Println("cannot find offer container")
		return nil, err
	}
	return s.FindOfferFromContainer(identifier, container)
}
