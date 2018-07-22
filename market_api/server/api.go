package server

import (
	"encoding/json"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"log"
	"net/http"
)

func (server *Server) handleRegistration(w http.ResponseWriter, r *http.Request) {
	register := &AccountRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(register)

	req := pb.CreateAccountRequest{
		Name:  register.Name,
		Email: register.Email,
	}

	res, err := server.rpcClient.DoCreateAccount(server.ctx, &req)
	if err != nil {
		log.Println("failed to make rpc call:" + err.Error())
	}
	log.Println("token:" + res.Token)

	account := &AccountResponse{
		Token: res.Token,
	}
	response, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (server *Server) handleAuthorisation(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) handleCreateHolding(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) handleCreateOffer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) handleAcceptOffer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) handleCloseOffer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
