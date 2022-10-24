package main

import (
	"fmt"
	"github.com/hosseintrz/suggestion_api/internal/db"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"os"
)

func main() {
	store := db.NewStore()
	s := server.NewServer(store)

	s.Router.POST("/signup", s.Signup)
	s.Router.POST("/login", s.Login)
	s.Router.POST("/suggestions", s.SubmitSuggestion)
	s.Router.GET("/suggestions", s.GetSuggestions)
	s.Router.GET("/", s.HealthCheck)

	port := "80"
	if val, ok := os.LookupEnv("HTTP_PORT"); ok {
		port = val
	}

	addr := "0.0.0.0:" + port
	errChan := s.Serve(addr)
	select {
	case err := <-errChan:
		fmt.Println(err.Error())
	}

}
