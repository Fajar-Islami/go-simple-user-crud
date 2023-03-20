package controller

import (
	"log"
	"net/http"

	"strconv"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/dtos"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type UsersController interface {
	GetUsersFetch(ctx echo.Context) error
	GetAllUsers(ctx echo.Context) error
	GetUserByID(ctx echo.Context) error
	CreateUsers(ctx echo.Context) error
	UpdateUsersByID(ctx echo.Context) error
	DeleteUsersByID(ctx echo.Context) error
}

type usersControllerImpl struct {
	usersusecase usecase.UsersUseCase
}

func NewUsersController(usersusecase usecase.UsersUseCase) UsersController {
	return &usersControllerImpl{
		usersusecase: usersusecase,
	}
}

// @Tags Users
// @Summary API Get Users Fetch
// @Router /user/fetch [get]
// @Param qs query dtos.FilterUsers true "Payload body [RAW]"
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) GetUsersFetch(ctx echo.Context) error {
	filter := new(dtos.FilterUsers)
	if err := ctx.Bind(filter); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.GetUsersFetch(ctx, *filter)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

// @Tags Users
// @Summary API Get Users
// @Router /user [get]
// @Param qs query dtos.FilterUsers true "Payload body [RAW]"
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) GetAllUsers(ctx echo.Context) error {
	filter := new(dtos.FilterUsers)
	if err := ctx.Bind(filter); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.GetAllUsers(ctx, *filter)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

// @Tags Users
// @Summary API Get Users
// @Router /user/{userid} [get]
// @Param   id   path  int  true  "User ID"
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) GetUserByID(ctx echo.Context) error {
	userid := ctx.Param("userid")
	usersIdInt, errConv := strconv.Atoi(userid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.GetUserByID(ctx, usersIdInt)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDGETDATA, "", res, http.StatusOK)
}

// @Tags Users
// @Summary API Get Users
// @Router /user [POST]
// @Param Content-Type header string true "content type request" Enums(application/json)
// @Param request body dtos.ReqCreateDataUser true "Payload body [RAW]"
// @Accept json
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) CreateUsers(ctx echo.Context) error {
	params := new(dtos.ReqCreateDataUser)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.CreateUsers(ctx, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDPOSTDATA, "", res, http.StatusOK)
}

// @Tags Users
// @Summary API Get Users
// @Router /user/{userid} [put]
// @Param   id   path  int  true  "User ID"
// @Param Content-Type header string true "content type request" Enums(application/json)
// @Param request body dtos.ReqUpdateDataUser true "Payload body [RAW]"
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) UpdateUsersByID(ctx echo.Context) error {
	userid := ctx.Param("userid")
	usersIdInt, errConv := strconv.Atoi(userid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	params := new(dtos.ReqUpdateDataUser)
	if err := ctx.Bind(params); err != nil {
		log.Println(err)
		return helper.BuildResponse(ctx, false, helper.FAILEDPOSTDATA, err.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.UpdateUsersByID(ctx, usersIdInt, *params)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDUPDATEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}

// @Tags Users
// @Summary API Delete Users
// @Router /user/{userid} [delete]
// @Param   id   path  int  true  "User ID"
// @Produces json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
func (ohco *usersControllerImpl) DeleteUsersByID(ctx echo.Context) error {
	userid := ctx.Param("userid")
	usersIdInt, errConv := strconv.Atoi(userid)
	if errConv != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDGETDATA, errConv.Error(), nil, http.StatusBadRequest)
	}

	res, err := ohco.usersusecase.DeleteUsersByID(ctx, usersIdInt)
	if err != nil {
		return helper.BuildResponse(ctx, false, helper.FAILEDDELETEDATA, err.Err.Error(), nil, err.Code)
	}

	return helper.BuildResponse(ctx, true, helper.SUCCEEDUPDATEDATA, "", res, http.StatusOK)
}
