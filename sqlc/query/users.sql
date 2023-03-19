-- name: GetUserByID :one
SELECT id,email,first_name,last_name,avatar,created_at,updated_at FROM users
WHERE id = ? and deleted_at is null LIMIT 1;

-- name: GetManyUser :many
SELECT id,email,first_name,last_name,avatar,created_at,updated_at FROM users
WHERE deleted_at is null and (email like ? or first_name like ? or last_name like ?)
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: GetCountManyUser :one
SELECT count(*) FROM users
WHERE deleted_at is null and (email like ? or first_name like ? or last_name like ?) LIMIT 1;

-- name: CreateUser :execlastid
INSERT INTO users (
  id,email,first_name,last_name,avatar
) VALUES ( ?,?,?,?,? );

-- name: UpdatePartialUsers :execlastid
UPDATE users 
SET 
    email = COALESCE(sqlc.arg(email), email),
    first_name = COALESCE(sqlc.arg(first_name), first_name),
    last_name = COALESCE(sqlc.arg(last_name), last_name),
    avatar = COALESCE(sqlc.arg(avatar), avatar),
    updated_at = now()
WHERE id = ? and deleted_at is null;

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = now()
WHERE id = ? and deleted_at is null;
