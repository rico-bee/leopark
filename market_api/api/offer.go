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
		log.Println("failed to find offers:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(offers)
	if err != nil {
		log.Println("failed to serialise assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("offers:" + string(data))
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

	log.Println("creating offer for asset:" + createOffer.Asset)

	createOfferReq := &pb.CreateOfferRequest{
		Label:       createOffer.Label,
		Description: createOffer.Description,
		Source:      mapMarketplaceHolding(createOffer.Source, createOffer.Asset, createOffer.SrcQuantity),
		Rules:       mapRules(createOffer.Rules),
		PrivateKey:  auth.PrivateKey,
	}

	if createOffer.Target != "" {
		log.Println("target is defined:" + createOffer.Target)
		createOfferReq.Target = mapMarketplaceHolding(createOffer.Target, createOffer.Asset, createOffer.TargetQuantity)
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

	newOffer := CreateOfferResponse{
		Id: res.Id,
	}
	data, _ := json.Marshal(newOffer)
	w.Write(data)
}

func mapOfferParticipant(srcId, targetId string, source, target *Holding) *pb.OfferParticipant {
	return &pb.OfferParticipant{
		SrcHolding:    srcId,
		TargetHolding: targetId,
		SrcAsset:      source.Asset,
		TargetAsset:   target.Asset,
	}
}

func (h *Handler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	auth, err := h.CurrentUser(w, r)
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Println("no id param specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acceptOffer := &AcceptOfferRequest{}
	bindRequestBody(r, acceptOffer)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	offer, err := h.Db.FindOffer(id)
	if err != nil {
		log.Println("failed to find offer from request " + id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ids := []string{acceptOffer.Source, acceptOffer.Target}

	holdings := h.Db.FetchHoldingsByIds(ids)
	if holdings == nil {
		log.Println("failed to find holdings from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acceptOfferReq := &pb.AcceptOfferRequest{
		Identifier: id,
		Sender:     mapOfferParticipant(offer.Source, offer.Target, holdings[acceptOffer.Source], holdings[acceptOffer.Target]),
		Receiver:   mapOfferParticipant(acceptOffer.Source, acceptOffer.Target, holdings[acceptOffer.Source], holdings[acceptOffer.Target]),
		PrivateKey: auth.PrivateKey,
	}
	res, err := h.RpcClient.DoAcceptOffer(ctx, acceptOfferReq)
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

func (h *Handler) CloseOffer(w http.ResponseWriter, r *http.Request) {
	auth, _ := h.CurrentUser(w, r)
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Println("no id param specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	closeOffer := &pb.CloseOfferRequest{
		Id:         id,
		PrivateKey: auth.PrivateKey,
	}

	res, err := h.RpcClient.DoCloseOffer(ctx, closeOffer)
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
