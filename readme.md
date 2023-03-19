# TODO
[ ] CRUD User
[ ] Get Detail Save and Get From redis
[ ] Connect to 3rd party
[ ] Update swagger
[ ] Deploy to fly.io
[ ] Add Github Action
[ ] Auto Deploy with GA
[ ] Migrate from redis to groupcache


    first_name = IF(CAST(@set_first_name AS bool) = true, @first_name, first_name),
    last_name = IF( CAST(sqlc.arg(set_last_name) AS bool)  = true, @last_name , last_name),
    avatar = IF( CAST(sqlc.arg(set_avatar) AS bool)  = true, CAST(@avatar AS TEXT), avatar),
    updated_at = now()


    last_name = IF( sqlc.arg(set_last_name)  = true, sqlc.arg(last_name), last_name),
    avatar = IF( sqlc.arg(set_avatar)  = true, sqlc.arg(avatar), avatar),

    WHERE deleted_at is null and (email like ? or first_name like ? or last_name like ?)