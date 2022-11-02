package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/suggestion_api/internal/config"
	"github.com/hosseintrz/suggestion_api/internal/db/cache"
	db "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/sirupsen/logrus"
)

type Server struct {
	DB     *db.Store
	Router *gin.Engine
	Cache  cache.Cache
}

func NewServer(conf *config.Config) (*Server, error) {
	store, err := db.NewStore(conf.DBConfig)
	c, err2 := cache.NewCache(conf.CacheConfig)
	if err != nil {
		logrus.Warn("db error")
		return nil, err
	}
	if err2 != nil {
		logrus.Warn("cache error")
		return nil, err2
	}

	return &Server{
		DB:     store,
		Router: gin.Default(),
		Cache:  c,
	}, nil
}

func (s *Server) Serve(addr string) chan error {
	errChan := make(chan error)
	go func() {
		errChan <- s.Router.Run(addr)
	}()
	return errChan
}
