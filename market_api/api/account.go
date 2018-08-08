package api

import (
	"encoding/json"
	crypto "github.com/rico-bee/leopark/crypto"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func (h *Handler) FindAccount(w http.ResponseWriter, r *http.Request) {
	req := &FindAccountRequest{}
	bindRequestBody(r, req)

	account, err := h.Db.FindUser(req.Email)

	if err != nil {
		log.Println("failed to find assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	res := &FindAccountResponse{
		Email:     account.Email,
		PublicKey: account.PublicKey,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Println("failed to serialise assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	register := &CreateAccountRequest{}
	bindRequestBody(r, register)
	req := pb.CreateAccountRequest{
		Name:     register.Name,
		Email:    register.Email,
		Password: register.Password,
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	res, err := h.RpcClient.DoCreateAccount(ctx, &req)
	if err != nil {
		log.Println("failed to make rpc call:" + err.Error())

	}
	hashPwd, err := crypto.HashPassword(register.Password)
	if err != nil {
		log.Println("cannot hash password:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	authInfo := &crypto.AuthInfo{
		Email:      register.Email,
		PwdHash:    hashPwd,
		PrivateKey: res.PrivateKey,
		PublicKey:  res.PublicKey,
	}

	err = h.Db.CreateUser(authInfo)
	if err != nil {
		log.Println("failed to create auth in db for " + authInfo.Email)
		w.WriteHeader(http.StatusInternalServerError)
	}

	tokenString, err := crypto.GenerateAuthToken(authInfo)
	if err != nil {
		log.Println("failed to create jwt token from auth:" + authInfo.Email)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var response []byte
	if err == nil && tokenString != "" {
		account := &AccountResponse{
			Token: tokenString,
		}
		response, _ = json.Marshal(account)
		w.WriteHeader(http.StatusOK)
	}
	w.Write(response)
}

func (h *Handler) FindAuthorisation(w http.ResponseWriter, r *http.Request) {
	authorise := &FindAuthorisationRequest{}
	bindRequestBody(r, authorise)
	// Contact the server and print out its response.

	auth, err := h.Db.FindUser(authorise.Email)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	if !crypto.CheckPasswordHash(authorise.Password, auth.PwdHash) {
		log.Println("invalid password....")
		w.WriteHeader(http.StatusUnauthorized)
	}
	tokenString, err := crypto.GenerateAuthToken(auth)
	if err != nil {
		log.Println("failed to create auth token:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	account := &AccountResponse{
		Token: tokenString,
	}
	response, _ := json.Marshal(account)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
