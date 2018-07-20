package server

import (
	"time"

	"github.com/labstack/echo"
)

//RequestContext for server context
type RequestContext struct {
	ID                        string
	requestHash               string
	requestAPIKeyHash         string
	auth                      bool
	timer                     time.Time
	httpResponseInternalError string
	errorMessage              string
}

//ErrorResponsePayload for http request
type ErrorResponsePayload struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

//ErrorResponse for http request
type ErrorResponse struct {
	Code       int
	Message    string
	StatusCode int
}

// createServerContext it will create
func (server *Server) createServerContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()
		if path == "/ping" || path == "/favicon.ico" {
			return next(c)
		}
		requestID := server.generateRequestID()
		requestContext := &RequestContext{
			ID:          requestID,
			timer:       time.Now(),
			requestHash: server.generateRequestHash(c.Request()),
		}

		c.Set(RequestContextName, requestContext)
		c.Set(RequestContextID, requestID)
		return next(c)
	}
}
