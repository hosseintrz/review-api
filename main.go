package main

import (
	"fmt"
	"github.com/hosseintrz/suggestion_api/internal/config"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"github.com/hosseintrz/suggestion_api/internal/server/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.GetConfig()
	logFatal(err)

	s, err := server.NewServer(conf)
	logFatal(err)
	routes.SetupRoutes(s, conf)

	logrus.Infof(conf.ServerAddress)
	errChan := s.Serve(conf.ServerAddress)
	select {
	case err := <-errChan:
		fmt.Println(err.Error())
	}

}

func logFatal(err error) {
	if err != nil {
		logrus.Fatalf(err.Error())
	}
}
