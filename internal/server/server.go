package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/suggestion_api/internal/db"
)

type Server struct {
	DB     *db.Store
	Router *gin.Engine
}

func NewServer(db *db.Store) *Server {
	return &Server{
		DB:     db,
		Router: gin.Default(),
	}
}

func (s *Server) Serve(addr string) chan error {
	errChan := make(chan error)
	go func() {
		errChan <- s.Router.Run(addr)
	}()
	return errChan
}
