-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, email, hashed_password)
VALUES (
   NOW(), NOW(), $1, $2
)
RETURNING id, email, updated_at, created_at, is_chirpy_red;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

-- name: GetUserByEmaill :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = $2,
    email = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeUserToRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
