package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof" //Profiling the API
	"os"
	"strconv"
	"time"
)

// Server context structure
type Server struct {
	version     string
	identity    string
	profiling   bool
	listenPort  int
	httpServer  *echo.Echo
	region      string
	environment string
	rpcClient   *pb.MarketClient
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
func NewServer(rpc *pb.MarketClient, version string, debug bool) (*Server, error) {

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	logger.Info("Starting " + serverIdentity + " - version " + version)
	if debug {
		logger.Level = logrus.DebugLevel
		logger.Info("Logging set to debug")
	}
	server := &Server{
		version:     version,
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

func (server *Server) initializeRouting() {
	e := echo.New()
	// e.Logger = echo.Logger{server.logger}
	e.Pre(server.createServerContext, middleware.CORSWithConfig(CORSPolicyConfig))
	// e.Use(server.loggerHook())
	//e.Use(server.setResponseHeaders)

	if server.profiling {
		dbg := e.Group("/debug")
		// expvar
		dbg.GET("/pprof/*", func(c echo.Context) error {
			w := c.Response().Writer
			r := c.Request()
			h, p := http.DefaultServeMux.Handler(r)
			if p != "" {
				h.ServeHTTP(w, r)
				return nil
			}
			return echo.NewHTTPError(http.StatusNotFound)
		})
	}

	e.GET("/ping", func(c echo.Context) error {
		c.Response().Header().Set(HeaderServer, serverHostHeader)
		if server.checkOnlineStatus() {
			return c.JSON(http.StatusOK, "ok")
		}
		return c.String(http.StatusBadRequest, "fail")
	})

	// e.GET("/healthcheck", func(c echo.Context) error {
	// 	return server.handleAPIHealthCheck(c)
	// }, server.checkIsMethodAllowed)

	// skylab := e.Group("/v0/skylab/")

	// if server.checkIfSecureConnectionValidationRequired() {
	// 	skylab.Use(server.checkIsSecureConnection)
	// }

	// skylab.GET("address/search", func(c echo.Context) error {
	// 	return server.handleAPIv1AddressSearch(c)
	// }, server.checkIsMethodAllowed)

	// skylab : spaceship membership management apis
	// skylab.POST("signup", func(c echo.Context) error {
	// 	return server.handleLeadSignupFromWebhook(c)
	// }, server.checkIsMethodAllowed)

	// skylab.POST("referral", func(c echo.Context) error {
	// 	return server.handleLeadFindOrCreateReferralCode(c)
	// }, server.checkIsMethodAllowed)
	server.logger.Info("Server listening on port ", strconv.Itoa(server.listenPort))
	server.httpServer = e
}

// Start starts the cache
func (server *Server) Start() {
	logrus.Println("starting server...")
	server.initializeRouting()
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(server.listenPort),
		ReadTimeout:  time.Duration(httpServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(httpServerWriteTimeout) * time.Second,
	}
	if err := server.httpServer.StartServer(httpServer); nil != err {
		server.logger.WithError(err).Errorln("Fail start server")
	}

	server.logger.Info("We are stopping")
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
