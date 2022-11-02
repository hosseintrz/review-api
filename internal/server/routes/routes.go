package routes

import (
	"github.com/hosseintrz/suggestion_api/internal/config"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"github.com/hosseintrz/suggestion_api/internal/server/handlers"
	"github.com/hosseintrz/suggestion_api/internal/server/middlewares"
)

func SetupRoutes(s *server.Server, conf *config.Config) {
	authHandler := handlers.NewAuthHandler(s, conf.AuthConfig)
	reviewHandler := handlers.NewReviewHandler(s)
	healthHandler := handlers.NewHealthHandler(s)
	movieHandler := handlers.NewMovieHandler(s)

	s.Router.POST("/signup", authHandler.Signup)
	s.Router.POST("/login", authHandler.Login)
	s.Router.POST("/refresh", authHandler.RefreshToken)

	// PROTECTED ROUTES
	apiProtected := s.Router.Group("")
	apiProtected.Use(middlewares.ValidateJWT(s, conf.AuthConfig))

	//reviews
	apiProtected.POST("/reviews", reviewHandler.SubmitReview)
	apiProtected.GET("/reviews", reviewHandler.GetReviews)
	//movies
	apiProtected.POST("/movies", movieHandler.CreateMovie)
	apiProtected.GET("/movies/:id", movieHandler.GetMovie)
	//health
	apiProtected.GET("/ping", healthHandler.HealthCheck)
}
