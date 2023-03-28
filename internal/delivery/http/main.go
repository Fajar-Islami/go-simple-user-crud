package http

import (
	"fmt"
	"net/http"
	"time"

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
	e.Validator = NewValidator()

	e.Any("", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(LoggerMiddleware(&cont.Logger.Log))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: func(c echo.Context) bool {
			requestPath := c.Request().URL.String()
			return requestPath == "/swagger/*"
		},
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			c.Logger().Error(err)
		},
		Timeout: time.Duration(cont.Apps.CtxTimeout * int(time.Second)),
	}))

	port := fmt.Sprintf("%s:%d", cont.Apps.Host, cont.Apps.HttpPort)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s/api/v1", cont.Apps.SwaggerAddress)
	utils.InitSnowflake()

	api := e.Group("/api/v1") // /api
	api.Any("", HealthCheck)
	api.Any("/health", HealthCheck)
	handler.UserHandler(api, cont, AuthMiddleware)

	e.Logger.Fatal(e.Start(port))
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
