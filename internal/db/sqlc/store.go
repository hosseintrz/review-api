package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hosseintrz/suggestion_api/internal/config"
	_ "github.com/lib/pq"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(conf config.DBConfig) (*Store, error) {
	db, err := sql.Open(conf.Driver, conf.Source)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:      db,
		Queries: New(db),
	}, nil
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("query err : %s , rollback err : %s", err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

type ReviewTxParams struct {
	Username string         `json:"username"`
	MovieID  int64          `json:"movie_id"`
	Rating   sql.NullString `json:"rating"`
	Text     sql.NullString `json:"text"`
}
type ReviewTxResult struct {
	ID       int64          `json:"id"`
	Username string         `json:"username"`
	MovieID  int64          `json:"movie_id"`
	Rating   sql.NullString `json:"rating"`
	Text     sql.NullString `json:"text"`
}

func (store *Store) SubmitReviewTx(ctx context.Context, params ReviewTxParams) (ReviewTxResult, error) {
	var result ReviewTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		review, err := queries.InsertReview(ctx, InsertReviewParams{
			Username: params.Username,
			MovieID:  params.MovieID,
			Rating:   params.Rating,
			Text:     params.Text,
		})
		if err != nil {
			return err
		}
		//update movie rating
		err = queries.UpdateMovieRating(ctx, params.MovieID)
		if err != nil {
			return err
		}

		result = ReviewTxResult{
			ID:       review.ID,
			Username: params.Username,
			MovieID:  review.MovieID,
			Rating:   review.Rating,
			Text:     review.Text,
		}
		return nil
	})
	return result, err
}
