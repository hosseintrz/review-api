package test

import (
	"context"
	"database/sql"
	"fmt"
	db2 "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestSubmitReviewTx(t *testing.T) {
	store, err := db2.NewStore(conf.DBConfig)
	require.NoError(t, err)
	users := make([]db2.User, 0)
	reviews := make([]db2.ReviewTxResult, 0)
	var movie db2.Movie

	setup := func(n int) {
		movie = CreateRandomMovie(t)
		for i := 0; i < n; i++ {
			user := CreateRandomUser(t)
			users = append(users, user)
		}
	}
	cleanup := func() {
		logrus.Info("entered cleanup")
		for _, review := range reviews {
			err := store.Queries.DeleteReview(context.Background(), review.ID)
			require.NoError(t, err)
		}
		err := store.Queries.DeleteMovie(context.Background(), movie.ID)
		require.NoError(t, err)
		for _, user := range users {
			err = store.Queries.DeleteUser(context.Background(), user.ID)
			require.NoError(t, err)
		}
	}

	setup(3)
	defer cleanup()

	ratings := make([]float32, 0)
	for _, user := range users {

		rating := util.RandomRating()
		ratings = append(ratings, rating)
		params := db2.ReviewTxParams{
			Username: user.Username,
			MovieID:  movie.ID,
			Text:     sql.NullString{String: util.RandomString(40), Valid: true},
			Rating:   sql.NullString{String: fmt.Sprintf("%.2f", rating), Valid: true},
		}

		review, err := store.SubmitReviewTx(context.Background(), params)
		reviews = append(reviews, review)

		//review tests
		require.NoError(t, err)

	}
	//movie test
	m, err := store.Queries.GetMovie(context.Background(), movie.ID)
	require.NoError(t, err)
	r, err := strconv.ParseFloat(m.Rating.String, 32)
	require.InDelta(t, avg(ratings...), r, 0.1)
}

func avg(nums ...float32) float32 {
	var sum float32
	for _, num := range nums {
		sum += num
	}
	return sum / float32(len(nums))
}
