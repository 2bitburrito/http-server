-- name: CreateChirp :one
INSERT INTO chirps (
    body, user_id, created_at, updated_at
) VALUES ( $1, $2, $3, $4)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetSingleChirp :one
SELECT * FROM chirps
WHERE id = $1;
