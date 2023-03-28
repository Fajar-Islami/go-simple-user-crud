package main

import (
	"flag"
	"log"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var rollback bool

func main() {
	flag.BoolVar(&rollback, "rollback", false, "")
	flag.Parse()

	container.Initcont(".env")
	cont := container.InitContainer("mysql")
	driver, err := mysql.WithInstance(cont.Mysqldb, &mysql.Config{})
	if err != nil {
		log.Println("err", err)
	}

	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Println("err", err)
	}

	log.Println("Running migration")

	if rollback {
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		log.Println("Rollback Done!!!")
	} else {
		if err := m.Up(); err != nil {
			log.Fatal("err migrate up ", err)
		}
		log.Println("Migration Done!!!")
	}
}
