package main

import (
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"

	rest "github.com/Fajar-Islami/go-simple-user-crud/internal/delivery/http"
)

func main() {
	contConf := container.InitContainer()
	defer contConf.Mysqldb.Close()

	rest.HTTPRouteInit(contConf)
}
