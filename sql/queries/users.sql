-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, email)
VALUES (
   NOW(), NOW(), $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

