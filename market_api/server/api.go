package server

import (
	"encoding/json"
	"fmt"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strings"
	"time"
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func bindRequestBody(r *http.Request, dto interface{}) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(dto)
}

func (server *Server) handleRegistration(w http.ResponseWriter, r *http.Request) {
	register := &AccountRequest{}
	bindRequestBody(r, register)
	req := pb.CreateAccountRequest{
		Name:     register.Name,
		Email:    register.Email,
		Password: register.Password,
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	res, err := server.rpcClient.DoCreateAccount(ctx, &req)
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

func (server *Server) handleAuthorisation(w http.ResponseWriter, r *http.Request) {
	authorise := &AuthoriseRequest{}
	bindRequestBody(r, authorise)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	authoriseReq := &pb.AuthoriseAccountRequest{Email: authorise.Email, Password: authorise.Password}
	res, err := server.rpcClient.DoAuthoriseAccount(ctx, authoriseReq)
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
