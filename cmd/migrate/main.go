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

	// cont := container.GetEnv()
	// var mysqlConfig mysqlclient.MysqlConf
	// err := cont.Unmarshal(&mysqlConfig)
	// if err != nil {
	// 	log.Println("err", err)
	// }

	// u := &url.URL{}
	// u.Scheme = "mysql"
	// u.User = url.UserPassword(mysqlConfig.Username, mysqlConfig.Password)
	// u.Host = mysqlConfig.Host + ":" + strconv.Itoa(mysqlConfig.Port)
	// u.Path = mysqlConfig.DbName
	// v := url.Values{}
	// v.Set("sslmode", "disable")
	// u.RawQuery = v.Encode()

	// dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?&charset=utf8mb4&parseTime=True&loc=Local", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DbName)
	// fmt.Println("dsn", dsn)

	// m, err := migrate.New("file://migrations", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
