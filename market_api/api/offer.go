package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func mapOfferStatus(s string) int32 {
	switch s {
	case "UNSET":
		return 0
	case "OPEN":
		return 1
	case "CLOSED":
		return 2
	}
	return -1
}

func mapOfferStatusToStr(s int32) string {
	switch s {
	case 0:
		return "UNSET"
	case 1:
		return "OPEN"
	case 2:
		return "CLOSED"
	}
	return ""
}

func (h *Handler) FindOffers(w http.ResponseWriter, r *http.Request) {
	_, err := h.CurrentUser(w, r)
	if err != nil {
		log.Println("failed to authenticate:" + err.Error())
		return
	}
	q := r.URL.Query()

	query := make(map[string]interface{})
	src := q.Get("src")
	if src != "" {
		query["source"] = src
	}
	target := q.Get("target")
	if target != "" {
		query["target"] = target
	}
	status := q.Get("status")
	if status != "" {
		query["status"] = status
	}

	offers, err := h.Db.FetchOffers(query)
	if err != nil {
		log.Println("failed to find assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(offers)
	if err != nil {
		log.Println("failed to serialise assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) FindOffer(w http.ResponseWriter, r *http.Request) {
	_, err := h.CurrentUser(w, r)
	if err != nil {
		log.Println("failed to authenticate:" + err.Error())
		return
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Println("no id param specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offer, err := h.Db.FindOffer(id)
	if err != nil {
		log.Println("no id param specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &FindOfferResponse{Offer: offer}
	data, err := json.Marshal(res)
	if err != nil {
		log.Println("cannot marshal response:" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println("warn:" + err.Error())
	}
}

func mapMarketplaceHolding(id, asset string, quantity int64) *pb.MarketplaceHolding {
	return &pb.MarketplaceHolding{
		HoldingId: id,
		Asset:     asset,
		Quantity:  quantity,
	}
}

func (h *Handler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	auth, err := h.CurrentUser(w, r)
	createOffer := &CreateOfferRequest{}
	bindRequestBody(r, createOffer)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	createOfferReq := &pb.CreateOfferRequest{
		Label:       createOffer.Label,
		Description: createOffer.Description,
		Source:      mapMarketplaceHolding(createOffer.Source, createOffer.Asset, createOffer.SrcQuantity),
		Target:      mapMarketplaceHolding(createOffer.Target, createOffer.Asset, createOffer.TargetQuantity),
		Rules:       mapRules(createOffer.Rules),
		PrivateKey:  auth.PrivateKey,
	}
	res, err := h.RpcClient.DoCreateOffer(ctx, createOfferReq)
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
