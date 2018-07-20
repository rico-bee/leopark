package server

import (
	"net/http"

	"github.com/labstack/echo/middleware"
)

var (
	corsAllowedDomains = []string{"http://localhost:3000", "https://www.devspaceship.com.au", "https://app.devspaceship.com.au", "https://web.devspaceship.com.au", "https://app.spaceship.com.au", "https://www.spaceship.com.au"}
	corsAllowedMethods = []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions}

	//CORSPolicyConfig config
	CORSPolicyConfig = middleware.CORSConfig{
		AllowOrigins: corsAllowedDomains,
		AllowMethods: corsAllowedMethods,
	}
)
