package test

import (
	"context"
	"database/sql"
	"fmt"
	db2 "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/util"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func CreateRandomMovie(t *testing.T) db2.Movie {
	params := db2.CreateMovieParams{
		Name: util.RandomString(10),
		Year: rand.Int31n(50) + 1970,
		Director: sql.NullString{
			String: fmt.Sprintf("%s %s", util.RandomString(5), util.RandomString(8)),
			Valid:  true,
		},
	}
	movie, err := testQueries.CreateMovie(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, movie)
	require.NotZero(t, movie.ID)

	return movie
}

func DeleteMovie(t *testing.T, id int64) {
	err := testQueries.DeleteMovie(context.Background(), id)
	require.NoError(t, err)

	m, err := testQueries.GetMovie(context.Background(), id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, m)
}

func TestCreateDeleteMovie(t *testing.T) {
	m := CreateRandomMovie(t)
	DeleteMovie(t, m.ID)
}
