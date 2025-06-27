-- name: GetUser :one
SELECT id, name, email, avatar_url FROM users WHERE id = $1;
