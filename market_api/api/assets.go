package api

import (
	"encoding/json"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func (h *Handler) FindAssets(w http.ResponseWriter, r *http.Request) {
	_, err := h.CurrentUser(w, r)
	if err != nil {
		log.Println("failed to authenticate:" + err.Error())
		return
	}
	log.Println("finding assets....")
	assets, err := h.Db.FindAssets()
	if err != nil {
		log.Println("failed to find assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(assets)
	if err != nil {
		log.Println("failed to serialise assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("found assets....")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) FindAsset(w http.ResponseWriter, r *http.Request) {
	_, err := h.CurrentUser(w, r)
	req := &FindAssetRequest{}
	bindRequestBody(r, req)
	asset, err := h.Db.FindAsset(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	res := &FindAssetResponse{Asset: asset}
	data, err := json.Marshal(res)
	if err != nil {
		log.Println("failed to serialise assets:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

func mapRules(rules []*Rule) []*pb.AssetRule {
	assetRules := []*pb.AssetRule{}
	for _, r := range rules {
		assetRules = append(assetRules, &pb.AssetRule{
			Type:  r.Type,
			Value: r.Value,
		})
	}
	return assetRules
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	auth, err := h.CurrentUser(w, r)
	createAsset := &CreateAssetRequest{}
	bindRequestBody(r, createAsset)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	createAssetReq := &pb.CreateAssetRequest{
		Name:        createAsset.Name,
		Description: createAsset.Description,
		Rules:       mapRules(createAsset.Rules),
		PrivateKey:  auth.PrivateKey,
	}
	res, err := h.RpcClient.DoCreateAsset(ctx, createAssetReq)
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
