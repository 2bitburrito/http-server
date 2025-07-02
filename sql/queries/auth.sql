-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens(
    token, updated_at, user_id, expires_at
) VALUES(
  $1, now(), $2, $3
  );

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
    WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = now(),
    updated_at = now()
WHERE token = $1;
