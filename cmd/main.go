package main

import (
	rest "github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
)

func main() {
	container.Initcont(".env")
	contConf := container.InitContainer()

	defer contConf.Mysqldb.Close()

	rest.HTTPRouteInit(contConf)
}
