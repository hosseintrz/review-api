-- name: InsertReview :one
INSERT INTO reviews(
    user_id, movie_id, rating, text
) VALUES(
    (
        SELECT id FROM users
        WHERE username=$1
    ),
    $2, $3, $4
)
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;