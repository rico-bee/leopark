package api

type FindAssetRequest struct {
	Name string `json:"name"`
}

type FindAssetResponse struct {
	Asset *Asset `json:"asset"`
}

type FindAccountRequest struct {
	PublicKey string `json:"public_key"`
}

type FindAccountResponse struct {
	Email     string     `json:"email"`
	PublicKey string     `json:"public_key"`
	Holdings  []*Holding `json:"holdings"`
}

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type FindAuthorisationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountResponse struct {
	Token string `json:"token"`
}

type CreateAssetRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rules       []*Rule `json:"rules"`
}

type FindOffersRequest struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Status string `json:"status"`
}

type FindOffersResponse struct {
	Offers []*Offer `json:"offers"`
}

type FindOfferResponse struct {
	Offer *Offer `json:"offer"`
}

type CreateOfferRequest struct {
	Asset          string  `json:"asset"`
	Label          string  `json:"label"`
	Description    string  `json:"description"`
	Source         string  `json:"source"`
	SrcQuantity    int64   `json:"src_quantity"`
	Target         string  `json:"target"`
	TargetQuantity int64   `json:"target_quantity"`
	Rules          []*Rule `json:"rules"`
}

type CreateOfferResponse struct {
	Id string `json:"id"`
}

type AcceptOfferRequest struct {
	ID     string `json:"id"`
	Count  int64  `json:"count"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type AcceptOfferResponse struct{}

type CreateHoldingRequest struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Asset       string `json:"asset"`
	Quantity    int64  `json:"quantity"`
}

type FetchHoldingRequest struct {
	Holdings []string `json:"holdings"`
}

type CreateHoldingResponse struct {
	Id string `json:"id"`
}
