-- name: CreateMovie :one
INSERT INTO movies(
    name, year, director, rating
) VALUES (
    $1, $2, $3, $4
         )
RETURNING *;

-- name: GetMovie :one
SELECT * FROM movies
WHERE id = $1 LIMIT 1;

-- name: ListMovies :many
SELECT * FROM movies;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;

-- name: UpdateMovieRating :exec
UPDATE movies
SET rating = (
    SELECT avg(rating) FROM reviews
    WHERE movie_id=$1
    )
WHERE id=$1;