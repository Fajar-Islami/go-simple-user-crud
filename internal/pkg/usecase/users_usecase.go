package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"math"
	"net/http"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/dtos"
	repositories "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/mysql"
	"github.com/go-redis/redis/v8"

	redisRepo "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/redis"
	reqresAPI "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/rest/reqres"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"

	"github.com/labstack/echo/v4"
)

type UsersUseCase interface {
	GetUsersFetch(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResDataUsers, err *helper.ErrorStruct)
	GetAllUsers(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResDataUsers, err *helper.ErrorStruct)
	GetUserByID(ctx echo.Context, userid int) (res dtos.ResDataUserSingle, err *helper.ErrorStruct)
	CreateUsers(ctx echo.Context, params dtos.ReqCreateDataUser) (res int, err *helper.ErrorStruct)
	UpdateUsersByID(ctx echo.Context, userid int, params dtos.ReqUpdateDataUser) (res string, err *helper.ErrorStruct)
	DeleteUsersByID(ctx echo.Context, userid int) (res string, err *helper.ErrorStruct)
}

type usersUseCaseImpl struct {
	userrepo  repositories.UsersRepo
	redis     redisRepo.RedisUsersRepository
	reqresAPI reqresAPI.ReqResAPI
}

func NewUsersUseCase(userrepo *repositories.UsersRepo, redis *redisRepo.RedisUsersRepository, reqresAPI *reqresAPI.ReqResAPI) UsersUseCase {
	return &usersUseCaseImpl{
		userrepo:  *userrepo,
		redis:     *redis,
		reqresAPI: *reqresAPI,
	}

}

func (usc *usersUseCaseImpl) GetUsersFetch(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResDataUsers, err *helper.ErrorStruct) {
	log := ctx.Logger()
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}
	dataRows := make([]dtos.ResDataUserSingle, 0)

	resAPI, errAPI := usc.reqresAPI.GetListUser(ctx.Request().Context(), reqresAPI.ReqListUser{
		PerPage: params.Limit,
		Page:    params.Page,
	})

	if errAPI != nil {
		log.Error("Err when get GetUsersFetch :", errAPI)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errAPI,
		}
	}

	usersData := make([]repositories.CreateUserParams, 0)

	for _, v := range resAPI.Data {
		dataRows = append(dataRows, dtos.ResDataUserSingle{
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar:    v.Avatar,
		})

		id := utils.IDGenerator()
		usersData = append(usersData, repositories.CreateUserParams{
			ID:        id,
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar: sql.NullString{
				String: v.Avatar,
				Valid:  true,
			},
		})
	}

	go func(usersData []repositories.CreateUserParams) {
		if len(usersData) < 1 {
			return
		}
		newCtx := context.Background()
		err := usc.userrepo.InsertBulkUsers(newCtx, usersData)
		if err != nil {
			log.Error(err)
		}

	}(usersData)

	res.Data = dataRows
	res.Page = resAPI.Page
	res.Rows = resAPI.PerPage
	res.TotalRows = resAPI.Total
	res.TotalPages = resAPI.TotalPages
	return res, nil
}

func (usc *usersUseCaseImpl) GetAllUsers(ctx echo.Context, params dtos.FilterUsers) (res dtos.ResDataUsers, err *helper.ErrorStruct) {
	log := ctx.Logger()

	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}
	dataRows := make([]dtos.ResDataUserSingle, 0)

	cpPage := params.Page
	params.Page = (params.Page - 1) * params.Limit

	paramsSearch := fmt.Sprintf("%%%s%%", params.Search)
	resRepo, errRepo := usc.userrepo.GetManyUser(ctx.Request().Context(), repositories.GetManyUserParams{
		Email:     paramsSearch,
		FirstName: paramsSearch,
		LastName:  paramsSearch,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Page),
	})

	if errors.Is(errRepo, sql.ErrNoRows) {
		log.Error("No data users error :", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("No Data Users"),
		}
	}

	if errRepo != nil {
		log.Error("Err when get many user :", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	countRows, errRepo := usc.userrepo.GetCountManyUser(ctx.Request().Context(), repositories.GetCountManyUserParams{
		Email:     paramsSearch,
		FirstName: paramsSearch,
		LastName:  paramsSearch,
	})

	if errRepo != nil {
		log.Error("Err when get count user :", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		dataRows = append(dataRows, dtos.ResDataUserSingle{
			DtosModel: dtos.DtosModel{
				ID:        &v.ID,
				CreatedAt: &v.CreatedAt.Time,
				UpdatedAt: &v.UpdatedAt.Time,
			},
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar:    v.Avatar.String,
		})
	}

	rows := params.Limit
	totalPage := math.Ceil(float64(countRows) / float64(rows))
	if rows > int(countRows) {
		rows = int(countRows)
	}

	res.Data = dataRows
	res.Page = cpPage
	res.Rows = rows
	res.TotalRows = int(countRows)
	res.TotalPages = int(totalPage)
	return res, nil
}

func (usc *usersUseCaseImpl) GetUserByID(ctx echo.Context, userid int) (res dtos.ResDataUserSingle, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	log := ctx.Logger()

	// Check data from redis
	data, errRedis := usc.redis.GetUsersCtx(contx, userid)
	if errRedis != nil {
		log.Error("Error when GetUsersCtx from redis: ", errRedis)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRedis,
		}
	}

	if data != nil {
		return *data, nil
	}

	// Check data from mysql
	resRepo, err := usc.getUserByIDHelper(contx, userid)
	if err != nil {
		log.Error("Error when getUserByIDHelper at GetUserByID: ", err)
		return res, err
	}
	res = dtos.ResDataUserSingle{
		DtosModel: dtos.DtosModel{
			ID:        &resRepo.ID,
			CreatedAt: &resRepo.CreatedAt.Time,
			UpdatedAt: &resRepo.UpdatedAt.Time,
		},
		Email:     resRepo.Email,
		FirstName: resRepo.FirstName,
		LastName:  resRepo.LastName,
		Avatar:    resRepo.Avatar.String,
	}

	// Set data to redis
	errRedis = usc.redis.SetUsersCtx(contx, &res)
	if errRedis != nil || errRedis == redis.Nil {
		log.Error("Error when SetUsersCtx from redis: ", errRedis)
	}

	return res, nil
}

func (usc *usersUseCaseImpl) CreateUsers(ctx echo.Context, params dtos.ReqCreateDataUser) (res int, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	err = usecaseValidation(ctx, params)
	log := ctx.Logger()
	if err != nil {
		return res, err
	}

	id := utils.IDGenerator()

	_, errRepo := usc.userrepo.CreateUser(contx, repositories.CreateUserParams{
		ID:        id,
		Email:     params.Email,
		FirstName: params.Firstname,
		LastName:  params.Lastname,
		Avatar: sql.NullString{
			String: params.Avatar,
			Valid:  true,
		},
	})

	if helper.MysqlCheckErrDuplicateEntry(errRepo) {
		log.Error("Err when CreateUser Duplicate :", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errors.New(errRepo.Error()),
		}
	}

	if errRepo != nil {
		log.Error("Err when CreateUser :", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return int(id), nil
}

func (usc *usersUseCaseImpl) UpdateUsersByID(ctx echo.Context, userid int, params dtos.ReqUpdateDataUser) (res string, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	err = usecaseValidation(ctx, params)
	if err != nil {
		return res, err
	}
	log := ctx.Logger()

	// Check data from mysql
	resRepo, err := usc.getUserByIDHelper(contx, userid)
	if err != nil {
		log.Error("Error when getUserByIDHelper at UpdateUsersByID: ", resRepo)
		return res, err
	}

	if params.Firstname == "" {
		params.Firstname = resRepo.FirstName
	}
	if params.Lastname == "" {
		params.Lastname = resRepo.LastName
	}
	if params.Email == "" {
		params.Email = resRepo.Email
	}

	avatarSql := sql.NullString{}
	if params.Avatar != "" {
		avatarSql.String = params.Avatar
		avatarSql.Valid = true
	}

	_, errRepo := usc.userrepo.UpdatePartialUsers(contx, repositories.UpdatePartialUsersParams{
		FirstName: params.Lastname,
		LastName:  params.Firstname,
		Avatar:    avatarSql,
		ID:        resRepo.ID,
		Email:     params.Email,
	})

	if errRepo != nil {
		log.Error("Error when UpdatePartialUsers : ", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	// // Delete key in redis
	errRedis := usc.redis.DeleteUsersCtx(ctx.Request().Context(), userid)
	if errRedis != nil {
		log.Error("Error when DeleteUsersCtx from redis: ", errRedis)
	}

	return "Succeed update user", nil
}

func (usc *usersUseCaseImpl) DeleteUsersByID(ctx echo.Context, userid int) (res string, err *helper.ErrorStruct) {
	contx := ctx.Request().Context()
	log := ctx.Logger()

	// Find user first
	// Check data from mysql
	resRepo, err := usc.getUserByIDHelper(contx, userid)
	if err != nil {
		log.Error("Error when getUserByIDHelper at DeleteUsersByID: ", err)
		return res, err
	}

	errRepo := usc.userrepo.SoftDeleteUser(contx, resRepo.ID)
	if errRepo != nil {
		log.Error("Error delete user : ", errRepo)
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	errRedis := usc.redis.DeleteUsersCtx(ctx.Request().Context(), userid)
	if errRedis != nil {
		log.Error("Error when DeleteUsersCtx from redis: ", errRedis)
	}

	return "Succeed delete user", nil
}

func (usc *usersUseCaseImpl) getUserByIDHelper(ctx context.Context, userid int) (res repositories.GetUserByIDRow, err *helper.ErrorStruct) {
	resRepo, errRepo := usc.userrepo.GetUserByID(ctx, int64(userid))
	if errors.Is(errRepo, sql.ErrNoRows) {
		return res, &helper.ErrorStruct{
			Code: http.StatusNotFound,
			Err:  errors.New("Data Users not found"),
		}
	}

	if errRepo != nil {
		return res, &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return resRepo, nil
}
