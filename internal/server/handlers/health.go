package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"net/http"
)

type HealthHandler struct {
	server *server.Server
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{server: s}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
