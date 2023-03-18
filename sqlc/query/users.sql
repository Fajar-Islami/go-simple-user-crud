-- name: GetUser :one
SELECT id,email,first_name,last_name,avatar,created_at,updated_at,deleted_at FROM users
WHERE id = ? and deleted_at is null LIMIT 1;

-- name: GetManyUser :many
SELECT id,email,first_name,last_name,avatar,created_at,updated_at,deleted_at FROM users
WHERE (email like $1 or first_name like $1 or last_name like $1) and deleted_at is null
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: GetCountManyUser :many
SELECT count(*) FROM users
WHERE (email like $1 or first_name like $1 or last_name like $1) and deleted_at is null LIMIT 1;

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
WHERE id = @id and deleted_at is null;

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = now()
WHERE id = ? and deleted_at is null;
