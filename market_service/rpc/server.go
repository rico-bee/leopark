package rpc

import (
	"encoding/hex"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	mktpb "github.com/rico-bee/leopark/market"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"github.com/rico-bee/leopark/market_service/transaction"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const (
	port = ":50051"
)

type server struct {
	ctx        signing.Context
	privateKey signing.PrivateKey
	signer     *signing.Signer
	validator  *transaction.SawtoothAPI
}

func (s *server) DoCreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	privateKey := s.ctx.NewRandomPrivateKey()
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	batches, signature := transaction.CreateAccount(signer, s.signer, in.Name, in.Email)
	if signature == "" {
		log.Fatal("Failed to create account")
	}
	err := s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	log.Println("creating account for : " + in.Email)
	batchIds := []string{signature}
	time.Sleep(5 * time.Second)
	committed, err := s.validator.CheckBatchStatus(batchIds)
	if !committed {
		return nil, err
	}
	return &pb.CreateAccountResponse{
		PrivateKey: privateKey.AsHex(),
		PublicKey:  signer.GetPublicKey().AsHex(),
	}, nil
}

func (s *server) DoCreateAsset(ctx context.Context, in *pb.CreateAssetRequest) (*pb.CreateAssetResponse, error) {
	pk, err := hex.DecodeString(in.PrivateKey)
	if err != nil {
		log.Println("failed to decode private key:" + err.Error())
		return nil, err
	}
	privateKey := signing.NewSecp256k1PrivateKey(pk)
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	rules := []*mktpb.Rule{}
	for _, rule := range in.Rules {
		rules = append(rules, &mktpb.Rule{
			Type:  mktpb.Rule_RuleType(rule.Type),
			Value: []byte(rule.Value),
		})
	}

	batches, signature := transaction.CreateAsset(signer, s.signer, in.Name, in.Description, rules)
	if signature == "" {
		log.Println("Failed to create asset")
	}
	err = s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	batchIds := []string{signature}
	time.Sleep(5 * time.Second)
	committed, err := s.validator.CheckBatchStatus(batchIds)
	if !committed {
		return nil, err
	}
	return &pb.CreateAssetResponse{Message: "success"}, nil
}

func (s *server) DoCreateHolding(ctx context.Context, req *pb.CreateHoldingRequest) (*pb.CreateHoldingResponse, error) {
	pk, err := hex.DecodeString(req.PrivateKey)
	if err != nil {
		log.Println("failed to decode private key:" + err.Error())
		return nil, err
	}
	privateKey := signing.NewSecp256k1PrivateKey(pk)
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	id, _ := uuid.GenerateUUID()
	batches, signature := transaction.CreateHolding(signer, s.signer,
		id,
		req.Label, req.Descrption, req.Asset, req.Quantity)

	if signature == "" {
		log.Println("Failed to create account")
	}
	err = s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	log.Println("holding created: " + req.Asset)
	batchIds := []string{signature}
	time.Sleep(5 * time.Second)
	committed, err := s.validator.CheckBatchStatus(batchIds)
	if !committed {
		return nil, err
	}
	return &pb.CreateHoldingResponse{Message: "sucess"}, nil
}

func (s *server) DoCreateOffer(ctx context.Context, req *pb.CreateOfferRequest) (*pb.CreateOfferResponse, error) {
	pk, err := hex.DecodeString(req.PrivateKey)
	if err != nil {
		log.Println("failed to decode private key:" + err.Error())
		return nil, err
	}
	privateKey := signing.NewSecp256k1PrivateKey(pk)
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	mktRules, err := MapAssetRule(req.Rules)
	if err != nil {
		return nil, err
	}
	id, _ := uuid.GenerateUUID()
	log.Println("we creating offer with src asset: " + req.Source.Asset)
	batches, signature := transaction.CreateOffer(signer, s.signer, id, req.Label, req.Description,
		MapHolding(req.Source), MapHolding(req.Target), mktRules)
	if signature == "" {
		log.Fatal("Failed to create account")
	}
	err = s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	time.Sleep(5 * time.Second)
	batchIds := []string{signature}
	committed, err := s.validator.CheckBatchStatus(batchIds)
	if !committed {
		return nil, err
	}
	return &pb.CreateOfferResponse{Message: "success"}, nil
}

func (s *server) DoAcceptOffer(ctx context.Context, req *pb.AcceptOfferRequest) (*pb.AcceptOfferResponse, error) {
	privateKey := signing.NewSecp256k1PrivateKey([]byte(req.PrivateKey))
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	batches, signature := transaction.AcceptOffer(signer, s.signer, req.Identifier,
		uint64(req.Count), MapOfferParticipant(req.Sender), MapOfferParticipant(req.Receiver))
	if signature == "" {
		log.Fatal("Failed to create account")
	}
	err := s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	return &pb.AcceptOfferResponse{Message: "success"}, nil
}

func (s *server) DoCloseOffer(ctx context.Context, req *pb.CloseOfferRequest) (*pb.CloseOfferResponse, error) {
	privateKey := signing.NewSecp256k1PrivateKey([]byte(req.PrivateKey))
	signer := signing.NewCryptoFactory(s.ctx).NewSigner(privateKey)
	batches, signature := transaction.CloseOffer(signer, s.signer, req.Id)
	if signature == "" {
		log.Fatal("Failed to create account")
	}
	err := s.validator.BatchRequest(batches)
	if err != nil {
		log.Println("failed to send batch request")
	}
	return &pb.CloseOfferResponse{Message: "success"}, nil
}

func newRpcServer() *server {
	api := transaction.NewSawtoothApi("tcp://localhost:4040")

	rpcServer := &server{
		ctx:       signing.CreateContext("secp256k1"),
		validator: api,
	}
	rpcServer.privateKey = signing.NewSecp256k1PrivateKey([]byte(RpcConfig.batchPrivateKey))
	rpcServer.signer = signing.NewCryptoFactory(rpcServer.ctx).NewSigner(rpcServer.privateKey)
	return rpcServer
}

//StartRpcServer starts RPC server
func StartRpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMarketServer(s, newRpcServer())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
