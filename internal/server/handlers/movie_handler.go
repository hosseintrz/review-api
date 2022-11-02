package handlers

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/internal/server"
	errors2 "github.com/hosseintrz/suggestion_api/internal/server/errors"
	"github.com/hosseintrz/suggestion_api/internal/server/responses"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	server *server.Server
}

func NewMovieHandler(s *server.Server) *MovieHandler {
	return &MovieHandler{
		server: s,
	}
}

type MovieDto struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Year     int32  `json:"year"`
	Director string `json:"director"`
	Rating   string `json:"rating"`
}

func NewMovieDto(m db.Movie) *MovieDto {
	return &MovieDto{
		ID:       m.ID,
		Name:     m.Name,
		Year:     m.Year,
		Director: m.Director.String,
		Rating:   m.Rating.String,
	}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var body MovieDto
	if err := c.ShouldBind(&body); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, errors2.ErrInvalidBody)
	}
	params := db.CreateMovieParams{
		Name:     body.Name,
		Year:     body.Year,
		Director: sql.NullString{String: body.Director, Valid: true},
		Rating:   sql.NullString{String: body.Rating, Valid: true},
	}
	movie, err := h.server.DB.CreateMovie(context.Background(), params)
	if err != nil {
		logrus.Warn("error creating movie -> ", err.Error())
		responses.ErrorResponse(c, http.StatusInternalServerError, errors2.ErrCreatingInstance)
	}
	movieDto := NewMovieDto(movie)
	responses.AbortResponse(c, http.StatusCreated, movieDto)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, errors2.ErrInvalidBody)
	}

	movie, err := h.server.DB.GetMovie(context.Background(), int64(id))
	if err != nil {
		responses.ErrorResponse(c, http.StatusNotFound, nil)
	}
	responses.AbortResponse(c, http.StatusOK, NewMovieDto(movie))
}
