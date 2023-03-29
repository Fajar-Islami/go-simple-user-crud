package handler

import (
	"testing"

	custommmiddleware "github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http/middleware"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
)

var contConf *container.Container

func TestMain(m *testing.M) {
	container.Initcont(".env.test")
	contConf = container.InitContainer()

	m.Run()
	contConf.Mysqldb.Close()
}

func TestInitContainer(t *testing.T) {
	container.Initcont(".env.test")
	contConf = container.InitContainer()

	defer contConf.Mysqldb.Close()

}

func HelperRouterUser(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Validator = custommmiddleware.NewValidator()
	// e.Use(middleware.Logger())
	utils.InitSnowflake()

	api := e.Group("/api/v1") // /api
	UserHandler(api, contConf, custommmiddleware.AuthMiddleware)
	return e
}
