package repositories

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

type UsersRepo interface {
	InsertBulkUsers(ctx context.Context, datas []CreateUserParams) (err error)
	CreateUser(ctx context.Context, arg CreateUserParams) (int64, error)
	GetCountManyUser(ctx context.Context, arg GetCountManyUserParams) (int64, error)
	GetManyUser(ctx context.Context, arg GetManyUserParams) ([]GetManyUserRow, error)
	GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error)
	UpdatePartialUsers(ctx context.Context, arg UpdatePartialUsersParams) (int64, error)
	SoftDeleteUser(ctx context.Context, id int64) error
}

type UsersRepoImpl struct {
	db *sql.DB
	*Queries
}

func NewUserRepo(db *sql.DB) UsersRepo {
	return &UsersRepoImpl{
		db:      db,
		Queries: New(db),
	}
}

func (ur *UsersRepoImpl) InsertBulkUsers(ctx context.Context, datas []CreateUserParams) (err error) {
	createUserSql := `-- name: CreateUserBulk
	INSERT INTO users (
	id,email,first_name,last_name,avatar
	) VALUES`
	vals := []interface{}{}

	for _, row := range datas {
		createUserSql += "( ?,?,?,?,? ),"
		vals = append(vals, row.ID, row.Email, row.FirstName, row.LastName, row.Avatar)
	}

	//trim the last ,
	createUserSql = strings.TrimSuffix(createUserSql, ",")

	//prepare the statement
	createUserSql += "ON DUPLICATE KEY UPDATE email=email;"
	stmt, err := ur.Queries.db.PrepareContext(ctx, createUserSql)
	if err != nil {
		return err
	}

	// format all vals at once
	res, err := stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}
	affectedRow, _ := res.RowsAffected()
	log.Println("createUserSql affectedRow : ", affectedRow)

	return nil
}
