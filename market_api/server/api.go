package server

import (
	"github.com/labstack/echo"
	"net/http"
)

func (server *Server) handleRegistration(c echo.Context) error {
	return nil
}

func (server *Server) handleAuthorisation(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func (server *Server) handleCreateAsset(c echo.Context) error {
	return nil
}

func (server *Server) handleCreateHolding(c echo.Context) error {
	return nil
}

func (server *Server) handleCreateOffer(c echo.Context) error {
	return nil
}

func (server *Server) handleAcceptOffer(c echo.Context) error {
	return nil
}

func (server *Server) handleCloseOffer(c echo.Context) error {
	return nil
}
