package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/internal/server"
	errors2 "github.com/hosseintrz/suggestion_api/internal/server/errors"
	"github.com/hosseintrz/suggestion_api/internal/server/responses"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrTokenGen         = errors.New("couldn't generate token")
	ErrReviewSubmission = errors.New("error while submitting review")
)

type ReviewHandler struct {
	server *server.Server
}

func NewReviewHandler(server *server.Server) *ReviewHandler {
	return &ReviewHandler{
		server: server,
	}
}

func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	userObj, _ := c.Get("user")
	user := userObj.(db.User)
	logrus.Info("username is ", user.Username)

	var body SubmitReviewReq
	err := c.ShouldBind(&body)
	if err != nil {
		logrus.Warn(err)
		responses.ErrorResponse(c, http.StatusBadRequest, errors2.ErrInvalidBody)
		return
	}

	res, err := h.server.DB.SubmitReviewTx(context.Background(), db.ReviewTxParams{
		Username: user.Username,
		MovieID:  body.MovieID,
		Text:     sql.NullString{String: body.Text, Valid: true},
		Rating:   sql.NullString{String: body.Rating, Valid: true},
	})
	if err != nil {
		logrus.Warn(ErrReviewSubmission.Error(), err.Error())
		responses.ErrorResponse(c, http.StatusInternalServerError, ErrReviewSubmission)
		return
	}

	c.JSON(201, gin.H{
		"ok": true,
		"id": res.ID,
	})
}
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	// TODO
	c.String(200, "todo")
}

type SubmitReviewReq struct {
	MovieID int64  `json:"movie_id"`
	Rating  string `json:"rating"`
	Text    string `json:"text"`
}
