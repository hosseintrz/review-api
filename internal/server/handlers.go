package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

type AuthReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (s *Server) Signup(c *gin.Context) {
	var req AuthReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "no username or password provided",
		})
		return
	}
	_, err := s.DB.GetUser(req.Username)
	if err == nil {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "user already exists",
		})
		return
	}
	token, err := s.DB.AddUser(req.Username, req.Password)
	if err != nil {
		c.JSON(500, "unexpected error")
		return
	}
	c.JSON(201, gin.H{
		"ok":    true,
		"token": token,
	})
}

func (s *Server) Login(c *gin.Context) {
	req := AuthReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "no username or password provided",
		})
		return
	}
	user, err := s.DB.GetUser(req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "invalid username or password",
		})
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"token": user.Token,
	})
}
func (s *Server) SubmitSuggestion(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.String(http.StatusUnauthorized, "empty authorization header")
		return
	}
	username, err := s.DB.GetUsernameFromToken(token)
	if err != nil {
		c.JSON(401, gin.H{
			"ok":    false,
			"error": "invalid token",
		})
		return
	}
	fmt.Println(username)
	text, ok := c.GetPostForm("text")
	if !ok || text == "" {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "no text provided",
		})
		return
	}
	s.DB.SubmitSuggest(username, text)
	c.JSON(201, gin.H{
		"ok": true,
	})
}
func (s *Server) GetSuggestions(c *gin.Context) {
	suggestions := s.DB.GetSuggestions()
	c.JSON(200, suggestions)
}
