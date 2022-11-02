package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/suggestion_api/internal/config"
	db "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"github.com/hosseintrz/suggestion_api/internal/server/responses"
	"github.com/hosseintrz/suggestion_api/internal/services"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("something failed during token generation")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthHandler struct {
	server       *server.Server
	tokenService *services.TokenService
	conf         config.AuthConfig
}

func NewAuthHandler(server *server.Server, conf config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		server:       server,
		tokenService: services.NewTokenService(server, conf),
		conf:         conf,
	}
}

type AuthReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req AuthReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "no username or password provided",
		})
		return
	}
	_, err := h.server.DB.GetUserByUsername(context.Background(), req.Username)
	if err == nil {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "user already exists",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user, err := h.server.DB.CreateUser(context.Background(), db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		logrus.Warn(err.Error())
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "error creating user",
		})
	}

	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, ErrTokenGeneration)
		return
	}

	response := NewTokenResponse(accessToken, refreshToken, exp)
	responses.Response(c, http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	req := AuthReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "no username or password provided",
		})
		return
	}
	user, err := h.server.DB.GetUserByUsername(context.Background(), req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidCredentials)
		return
	}
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, ErrTokenGeneration)
		return
	}
	response := NewTokenResponse(accessToken, refreshToken, exp)
	responses.Response(c, http.StatusOK, response)
}

type TokenRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	req := TokenRequest{}
	if err := c.ShouldBind(&req); err != nil {
		//		c.AbortWithError(http.StatusBadRequest, err)
		responses.AbortResponse(c, http.StatusBadRequest, nil)

	}
	claims, err := h.tokenService.ParseToken(req.Token, h.conf.RefreshSecret)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidToken)
		return
	}
	user, err := h.tokenService.ValidateToken(claims, true)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, ErrInvalidToken)
		return
	}
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, ErrTokenGeneration)
		return
	}
	response := NewTokenResponse(
		accessToken, refreshToken, exp,
	)

	responses.Response(c, http.StatusOK, response)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

func NewTokenResponse(access, refresh string, exp int64) *TokenResponse {
	return &TokenResponse{access, refresh, exp}
}
