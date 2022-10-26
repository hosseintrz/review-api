-- name: CreateSuggestion :one
INSERT INTO suggestions(
    user_id, text
) VALUES(
    (
        SELECT id FROM users
        WHERE username=$1
    ),
    $2
)
RETURNING *;