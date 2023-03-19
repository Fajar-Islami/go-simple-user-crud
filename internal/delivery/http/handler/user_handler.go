package handler

import (
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/controller"
	mysqlrepo "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/mysql"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/usecase"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthHandler(r *echo.Group, containerConf *container.Container, authMiddleware utils.MiddlewareType) {
	repo := mysqlrepo.New(containerConf.Mysqldb)
	usecase := usecase.NewUsersUseCase(*repo)
	controller := controller.NewUsersController(usecase)

	usersAPI := r.Group("/user")

	usersAPI.GET("", controller.GetAllUsers)
	usersAPI.GET("/:userid", controller.GetUserByID)
	usersAPI.POST("", controller.CreateUsers)
	usersAPI.PUT("/:userid", controller.UpdateUsersByID)
	usersAPI.DELETE("/:userid", authMiddleware(controller.DeleteUsersByID))

}
