package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof" //Profiling the API
	"os"
)

// Server context structure
type Server struct {
	identity    string
	profiling   bool
	listenPort  int
	region      string
	environment string
	rpcClient   pb.MarketClient
	logger      *logrus.Logger
}

const (
	version                   string = "0.0.1"
	region                    string = "sydney"
	serverIdentity            string = "api.market"
	serverHostHeader          string = "leopark-api"
	httpMaxIdleConnections    int    = 30
	httpRequestTimeout        int    = 60
	httpServerReadTimeout     int    = 61
	httpServerWriteTimeout    int    = 120
	httpStatusUnauthorizedTTL int    = 300
	httpServerOfflineFilePath string = "/tmp/api.leopark.offline.lock"
	environmentNoSecureCheck  string = "local"
)

// NewServer Server
func NewServer(rpc pb.MarketClient) (*Server, error) {

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	logger.Info("Starting " + serverIdentity + " - version " + version)
	server := &Server{
		identity:    serverIdentity,
		region:      region,
		listenPort:  8025,
		profiling:   false,
		environment: "dev",
		logger:      logger,
		rpcClient:   rpc,
	}
	return server, nil
}

// Start starts the cache
func (server *Server) Start() {
	logrus.Println("starting server...")
	r := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r.HandleFunc("/account", server.handleRegistration).Methods("POST")
	r.HandleFunc("/authorise", server.handleAuthorisation).Methods("POST")
	http.ListenAndServe(":8088", handlers.CORS(originsOk, headersOk, methodsOk)(r))
	//Stop Events go here
	server.logger.Info("We stopped successfully")
}

func (server *Server) checkOnlineStatus() bool {
	if _, err := os.Stat(httpServerOfflineFilePath); err == nil {
		return false
	}
	return true
}

func (server *Server) checkIfSecureConnectionValidationRequired() bool {
	if server.environment == environmentNoSecureCheck {
		return false
	}
	return true
}
