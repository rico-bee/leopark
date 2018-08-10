package server

import (
	//"github.com/auth0/go-jwt-middleware"
	//jwt "github.com/dgrijalva/jwt-go"//
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// crypto "github.com/rico-bee/leopark/crypto"
	api "github.com/rico-bee/leopark/market_api/api"
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
	logger      *logrus.Logger
	api         *api.Handler
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
func NewServer(handler *api.Handler) (*Server, error) {

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
		api:         handler,
	}
	return server, nil
}

// Start starts the cache
func (server *Server) Start() {
	logrus.Println("starting server...")
	r := mux.NewRouter()
	m := r.PathPrefix("/market").Subrouter()
	// !!! WARN: we have to explicitly whitelist all required headers in X-Requested-With to allow the pre-flight request pass
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// no jwt check on register
	r.HandleFunc("/register", server.api.CreateAccount).Methods("POST")
	r.HandleFunc("/authorise", server.api.FindAuthorisation).Methods("POST")

	m.HandleFunc("/account", server.api.FindAccount).Methods("GET")
	m.HandleFunc("/asset", server.api.CreateAsset).Methods("POST")
	m.HandleFunc("/asset/list", server.api.FindAssets).Methods("GET")
	m.HandleFunc("/asset/{name}", server.api.FindAsset).Methods("GET")
	m.Use(jwtMiddleware)
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	http.ListenAndServe(":8088", corsHandler)
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
