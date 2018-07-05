package rpc

import (
	pb "bitbucket.org/riczha/marketplace/market_service/proto/api"
	"bitbucket.org/riczha/marketplace/market_service/transaction"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	ctx       signing.Context
	signer    *signing.Signer
	validator *transaction.SawtoothAPI
}

func (s *server) DoCreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	//publicKey := signer.GetPublicKey().AsHex() // we only need it to create auth token
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
	return &pb.CreateAccountResponse{Message: "Hello " + in.Name}, nil
}

func newRpcServer() *server {
	api := transaction.NewSawtoothApi("tcp://localhost:4040")
	rpcServer := &server{
		ctx:       signing.CreateContext("secp256k1"),
		validator: api,
	}
	privateKey := signing.NewSecp256k1PrivateKey([]byte(RpcConfig.batchPrivateKey))
	rpcServer.signer = signing.NewCryptoFactory(rpcServer.ctx).NewSigner(privateKey)
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
