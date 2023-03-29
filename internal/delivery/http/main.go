package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/docs"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http/handler"
	custommmiddleware "github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http/middleware"
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
	e.Validator = custommmiddleware.NewValidator()

	e.Any("", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(custommmiddleware.LoggerMiddleware(&cont.Logger.Log))
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
	handler.UserHandler(api, cont, custommmiddleware.AuthMiddleware)

	// Start server
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server : ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("shutting down the server :", err)
	}
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
