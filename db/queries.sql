-- name: UserAdd :exec
INSERT INTO users (username, hashed_password, created_on, last_login)
VALUES ($1, $2, $3, $4);

-- name: UserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: SessionAdd :exec
INSERT INTO sessions (user_id, token, created_on, expires_on)
VALUES ($1, $2, $3, $4);

-- name: Session :one
SELECT user_id, token
FROM sessions
WHERE id = $1
AND expires_on > NOW();

