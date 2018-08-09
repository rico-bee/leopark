package api

import (
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type Handler struct {
	RpcClient pb.MarketClient
	Db        *DbServer
}

func NewHandler(rpc pb.MarketClient, db *DbServer) *Handler {
	return &Handler{RpcClient: rpc, Db: db}
}

func bindRequestBody(r *http.Request, dto interface{}) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(dto)
}

type rpcFunc func(context.Context, interface{}, ...interface{}) (interface{}, error)

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

// FromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func checkJwt(r *http.Request) (string, error) {
	tokenStr, err := FromAuthHeader(r)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) (string, error) {
	privateKey := r.Header.Get("privateKey")
	if privateKey == "" {
		unauthorised(w)
	}
	return privateKey, nil
}

func unauthorised(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorised access"))
}
