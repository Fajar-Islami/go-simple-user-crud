package http

import (
	"fmt"

	"github.com/Fajar-Islami/go-simple-user-crud/docs"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http/handler"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Simple User CRUD
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func HTTPRouteInit(cont *container.Container) {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Validator = NewValidator()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(LoggerMiddleware(&cont.Logger.Log))

	port := fmt.Sprintf("%s:%d", cont.Apps.Host, cont.Apps.HttpPort)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s/api/v1", port)
	utils.InitSnowflake()

	api := e.Group("/api/v1") // /api
	handler.AuthHandler(api, cont, AuthMiddleware)

	e.Logger.Fatal(e.Start(port))
}