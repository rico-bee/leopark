package server

type AccountRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AccountResponse struct {
	Token string `json:"token"`
}
