-- name: GetUser :one
SELECT id,email,first_name,last_name,avatar,created_at,updated_at,deleted_at FROM users
WHERE id = ? LIMIT 1;

-- name: GetManyUser :many
SELECT id,email,first_name,last_name,avatar,created_at,updated_at,deleted_at FROM users
WHERE email = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: GetCountManyUser :many
SELECT count(*) FROM users
WHERE email = ?;

-- name: CreateAuthor :execlastid
INSERT INTO users (
  id,email,first_name,last_name,avatar
) VALUES ( ?,?,?,?,? );

-- name: UpdatePartialUsers :execlastid
UPDATE users 
SET 
    first_name = IF(@update_first_name = true, @first_name, first_name),
    last_name = IF(@update_last_name = true, @last_name, last_name),
    avatar = IF(@update_avatar = true, @avatar, avatar),
    updated_at = now()
WHERE id = @id;

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = now()
WHERE id = ?;
