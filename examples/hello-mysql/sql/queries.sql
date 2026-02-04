-- name: GetUser :one
SELECT id, name, email
FROM users
WHERE id = ?
LIMIT 1;
