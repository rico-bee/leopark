package server

type AccountRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AccountResponse struct {
	Token string `json:"token"`
}

type AuthoriseRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAssetRequest struct {
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Rules       []*Rule `json: "rules"`
}
