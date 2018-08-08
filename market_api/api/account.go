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
	assets, err := h.Db.FindAssets()
	if err != nil {
		log.Println("failed to find assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	data, err := json.Marshal(assets)
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
		return
	}
	account := &AccountResponse{
		Token: res.Token,
	}
	response, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *Handler) FindAuthorisation(w http.ResponseWriter, r *http.Request) {
	authorise := &FindAuthorisationRequest{}
	bindRequestBody(r, authorise)
	// Contact the server and print out its response.

	auth, err := h.Db.FindUser(authorise.Email)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
