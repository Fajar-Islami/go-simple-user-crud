package http

import (
	"os"
	"testing"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http/handler"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var contConf *container.Container

func TestMain(m *testing.M) {
	container.Initcont(".env.test")
	contConf = container.InitContainer()

	defer contConf.Mysqldb.Close()

	os.Exit(m.Run())
}

func TestInitContainer(t *testing.T) {
	container.Initcont(".env.test")
	contConf = container.InitContainer()

	defer contConf.Mysqldb.Close()

}

func HelperRouterUser(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Validator = NewValidator()
	e.Use(middleware.Logger())

	api := e.Group("/api/v1") // /api
	handler.UserHandler(api, contConf, AuthMiddleware)
	return e
}
