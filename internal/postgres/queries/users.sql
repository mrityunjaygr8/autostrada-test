-- name: UserRetrieve :one
SELECT email, created,  id, admin FROM users WHERE id = $1 LIMIT 1;

-- name: UserRetrieveByEmail :one
SELECT email, created,  id, admin FROM users WHERE email = $1 LIMIT 1;

-- name: UsersList :many
WITH row_data AS (
    SELECT email, created, id, admin FROM users ORDER BY email LIMIT $1 OFFSET $2
) SELECT
      *,
      (SELECT COUNT(*) FROM users) AS row_data
FROM row_data;

-- name: UserInsert :one
INSERT INTO users (email, hashed_password, id, admin) VALUES ($1, $2, $3, $4) RETURNING email, created, id, admin;

-- name: UserUpdatePassword :execresult
UPDATE users SET hashed_password = $2 WHERE id = $1;

-- name: UserUpdateAdmin :execresult
UPDATE users SET admin = $2 WHERE id = $1;

-- name: UserDelete :execresult
DELETE FROM users WHERE id = $1;