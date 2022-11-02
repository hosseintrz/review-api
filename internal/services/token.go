package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hosseintrz/suggestion_api/internal/config"
	db "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/internal/server"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrCouldNotParseClaims = errors.New("couldn't parse claims")
)

const (
	ExpireAccessMin  = 30
	ExpireRefreshMin = 120
	AutoLogoutMin    = 10
)

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	ID  int64  `json:"id"`
	UID string `json:"uid"`
}

type TokenService struct {
	server *server.Server
	conf   config.AuthConfig
}

func NewTokenService(server *server.Server, conf config.AuthConfig) *TokenService {
	return &TokenService{
		server: server,
		conf:   conf,
	}
}

func (s *TokenService) GenerateToken(userID int64, expMin int, secret string) (string, string, int64, error) {
	exp := time.Now().Add(time.Minute * time.Duration(expMin))
	claims := &JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: exp}},
		ID:               userID,
		UID:              uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", 0, err
	}
	return signed, claims.UID, exp.Unix(), nil
}

type CachedToken struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

func (s *TokenService) GenerateTokenPair(user db.User) (
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID string
	accessToken, accessUID, exp, err = s.GenerateToken(user.ID, ExpireAccessMin, s.conf.AccessSecret)
	if err != nil {
		return
	}

	refreshToken, refreshUID, _, err = s.GenerateToken(user.ID, ExpireRefreshMin, s.conf.RefreshSecret)
	if err != nil {
		return
	}

	record, err := json.Marshal(CachedToken{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})
	if err != nil {
		return
	}
	err = s.server.Cache.Set(context.Background(), fmt.Sprintf("token-%d", user.ID), record, time.Minute*AutoLogoutMin)
	return
}

func (s *TokenService) ValidateToken(claims *JwtCustomClaims, isRefresh bool) (
	user db.User,
	err error,
) {
	var g errgroup.Group
	g.Go(func() error {
		token, err := s.server.Cache.Get(context.Background(), fmt.Sprintf("token-%d", claims.ID))
		if err != nil {
			return err
		}
		var cachedToken CachedToken
		err = json.Unmarshal([]byte(token), &cachedToken)

		var uid string
		if isRefresh {
			uid = cachedToken.RefreshUID
		} else {
			uid = cachedToken.AccessUID
		}

		if err != nil || uid != claims.UID {
			return ErrTokenNotFound
		}
		return nil
	})

	g.Go(func() error {
		user, err = s.server.DB.GetUser(context.Background(), claims.ID)
		if err != nil {
			return err
		}
		return nil
	})
	err = g.Wait()
	return
}

func (s *TokenService) ParseToken(tokenString, secret string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//	return nil, fmt.Errorf("unexpected signing method %s ", token.Header["alg"])
		//}
		return []byte(secret), nil
	})

	if err != nil {
		logrus.Warn("error ParseWithClaims -> ", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrCouldNotParseClaims
}
