package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/suggestion_api/internal/config"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"github.com/hosseintrz/suggestion_api/internal/server/responses"
	"github.com/hosseintrz/suggestion_api/internal/services"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidAuth = errors.New("Invalid Authorization")
)

func ValidateJWT(s *server.Server, conf config.AuthConfig) gin.HandlerFunc {
	tokenService := services.NewTokenService(s, conf)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn(ErrInvalidAuth.Error())
			responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuth)
			return
		}
		parts := strings.Split(authHeader, " ")
		authType, tokenStr := parts[0], parts[1]
		if authType != "JWT" {
			logrus.Warn(ErrInvalidAuth)
			responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuth)
			return
		}
		claims, err := tokenService.ParseToken(tokenStr, conf.AccessSecret)
		if err != nil {
			logrus.Warn("error parsing token -> ", err.Error())
			responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuth)
			return
		}
		user, err := tokenService.ValidateToken(claims, false)
		if err != nil {
			logrus.Warn("md : token is invalid -> ", err.Error())
			responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuth)
			return
		}
		c.Set("user", user)

		go func() {
			_ = s.Cache.Expire(context.Background(),
				fmt.Sprintf("token-%d", claims.ID),
				time.Minute*services.AutoLogoutMin)
		}()

		c.Next()
	}
}
