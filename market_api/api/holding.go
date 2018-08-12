package api

import (
	"encoding/json"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func PrintObj(msg interface{}) string {
	data, _ := json.Marshal(msg)
	return string(data)
}

func (h *Handler) CreateHolding(w http.ResponseWriter, r *http.Request) {
	auth, err := h.CurrentUser(w, r)
	createholding := &CreateHoldingRequest{}
	bindRequestBody(r, createholding)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	log.Println("creating holding:" + PrintObj(createholding))
	createHoldingReq := &pb.CreateHoldingRequest{
		Label:      createholding.Label,
		Descrption: createholding.Description, /// fix the typo todo
		Asset:      createholding.Asset,
		Quantity:   createholding.Quantity,
		PrivateKey: auth.PrivateKey,
	}
	res, err := h.RpcClient.DoCreateHolding(ctx, createHoldingReq)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to make rpc call:" + err.Error())
		log.Println("holding created failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Println("holding created ok")
		w.WriteHeader(http.StatusOK)
	}
	newHolding := &CreateHoldingResponse{Id: res.Id}
	data, _ := json.Marshal(newHolding)
	w.Write(data)
}
