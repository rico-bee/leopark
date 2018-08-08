package api

type FindAssetRequest struct {
	Name string `json: "name"`
}

type FindAssetResponse struct {
	Asset *Asset `json: "asset"`
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
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Rules       []*Rule `json: "rules"`
}
