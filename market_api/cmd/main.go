package main

import (
	// uuid "github.com/hashicorp/go-uuid"
	crypto "github.com/rico-bee/leopark/crypto"
	mktpb "github.com/rico-bee/leopark/market"
	api "github.com/rico-bee/leopark/market_api/api"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

const (
	rpcUrl = "localhost:50051"
)

var (
	app     = kingpin.New("leopark-cmd", "A command-line data import tool.")
	dataUri = app.Flag("file", "file path to load market data").String()
)

type Rule struct {
	RuleType string `yaml:"type"`
	Value    string `yaml:"value", omitempty`
}

type Asset struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Rules       []*Rule `yaml:"rules"`
}

type Holding struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Asset       string `yaml:"asset"`
	Quantity    int64  `yaml:"quantity"`
}

type Offer struct {
	Label          string  `yaml:"label"`
	Description    string  `yaml:"description"`
	Source         string  `yaml:"source"`
	SourceQuantity int64   `yaml:"sourceQuantity"`
	Rules          []*Rule `yaml:"rules"`
}

type Participant struct {
	Label       string       `yaml:"label"`
	Description string       `yaml:"description"`
	Email       string       `yaml:"email"`
	Password    string       `yaml:"password"`
	Assets      []*Asset     `yaml:"assets"`
	Holdings    []*Holding   `yaml:"holdings"`
	Offers      []*Offer     `yaml:"offers"`
	Renewables  []*Renewable `yaml:"renewables"`
}

type Renewable struct {
	Label          string
	Description    string
	Source         string
	SourceQuantity int64
	Rules          []Rule
}

type MarketData struct {
	Participants []*Participant `yaml:"accounts"`
}

func mapRules(rules []*Rule) []*pb.AssetRule {
	assetRules := []*pb.AssetRule{}
	for _, r := range rules {
		assetRule := &pb.AssetRule{
			Type:  mktpb.Rule_RuleType_value[r.RuleType],
			Value: r.Value,
		}
		assetRules = append(assetRules, assetRule)
	}
	return assetRules
}

func main() {

	data, err := ioutil.ReadFile("./app_data.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	mktData := MarketData{}
	err = yaml.Unmarshal(data, &mktData)
	if err != nil {
		log.Fatalln(err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(rpcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketClient(conn)

	db, err := api.NewDBServer("localhost:28015")
	if err != nil {
		log.Fatal("failed to connect to DB")
	}

	for _, p := range mktData.Participants {

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		req := &pb.CreateAccountRequest{
			Name:     p.Label,
			Email:    p.Email,
			Password: p.Password,
		}
		res, err := c.DoCreateAccount(ctx, req)
		if err != nil {
			log.Fatal("rpc failed:" + err.Error())
		}

		hashPwd, err := crypto.HashPassword(p.Password)
		if err != nil {
			log.Println("cannot hash password:" + err.Error())
		}
		auth := &crypto.AuthInfo{
			PublicKey:  res.PublicKey,
			PrivateKey: res.PrivateKey,
			PwdHash:    hashPwd,
			Email:      p.Email,
		}

		err = db.CreateUser(auth)
		if err != nil {
			log.Println("failed to create auth for user:" + auth.Email)
		}

		for _, a := range p.Assets {
			aReq := &pb.CreateAssetRequest{
				Name:        a.Name,
				Description: a.Description,
				Rules:       mapRules(a.Rules),
				PrivateKey:  res.PrivateKey,
			}
			_, err := c.DoCreateAsset(ctx, aReq)
			if err != nil {
				log.Fatal("rpc failed:" + err.Error())
			}
		}

		for _, h := range p.Holdings {
			if err != nil {
				log.Fatal("failed toget uuid:" + err.Error())
			}
			hReq := &pb.CreateHoldingRequest{
				Label:      h.Label,
				Descrption: h.Description,
				Asset:      h.Asset,
				Quantity:   h.Quantity,
				PrivateKey: res.PrivateKey,
			}
			_, err = c.DoCreateHolding(ctx, hReq)
			if err != nil {
				log.Println("rpc failed:" + err.Error())
			}
		}
	}
}
