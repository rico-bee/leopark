package api

import (
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func (h *Handler) CreateHolding(w http.ResponseWriter, r *http.Request) {
	auth, err := h.CurrentUser(w, r)
	createholding := &CreateHoldingRequest{}
	bindRequestBody(r, createholding)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	createAssetReq := &pb.CreateHoldingRequest{
		Label:      createholding.Label,
		Descrption: createholding.Description, /// fix the typo todo
		Asset:      createholding.Asset,
		Quantity:   createholding.Quantity,
		PrivateKey: auth.PrivateKey,
	}
	res, err := h.RpcClient.DoCreateHolding(ctx, createAssetReq)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to make rpc call:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write([]byte(res.Message))
}
