package handler

import (
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/controller"
	mysqlrepo "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/mysql"
	redisRepo "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/redis"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/rest"
	reqresAPI "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/rest/reqres"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/usecase"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthHandler(r *echo.Group, containerConf *container.Container, authMiddleware utils.MiddlewareType) {
	reqresRepo := reqresAPI.NewReqResAPI(containerConf.ReqResAPI.URL, rest.Opts{
		Timeout:     time.Duration(containerConf.ReqResAPI.TimeOut * int(time.Minute)),
		Logger:      *containerConf.Logger,
		IsDebugging: containerConf.ReqResAPI.Debugging,
	})
	redis := redisRepo.NewRedisRepoUsers(containerConf.Redis, &containerConf.Logger.Log)
	repo := mysqlrepo.NewUserRepo(containerConf.Mysqldb)
	usecase := usecase.NewUsersUseCase(&repo, &redis, &reqresRepo)
	controller := controller.NewUsersController(usecase)

	usersAPI := r.Group("/user")

	usersAPI.GET("/fetch", controller.GetUsersFetch)
	usersAPI.GET("", controller.GetAllUsers)
	usersAPI.GET("/:userid", controller.GetUserByID)
	usersAPI.POST("", controller.CreateUsers)
	usersAPI.PUT("/:userid", controller.UpdateUsersByID)
	usersAPI.DELETE("/:userid", authMiddleware(controller.DeleteUsersByID))

}
