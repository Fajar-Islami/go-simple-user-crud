package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"

	"github.com/spf13/viper"
)

type MysqlConf struct {
	Username           string `mapstructure:"mysql_username"`
	Password           string `mapstructure:"mysql_password"`
	DbName             string `mapstructure:"mysql_Dbname"`
	Host               string `mapstructure:"mysql_host"`
	Port               int    `mapstructure:"mysql_port"`
	Schema             string `mapstructure:"mysql_schema"`
	LogMode            bool   `mapstructure:"mysql_logMode"`
	MaxLifetime        int    `mapstructure:"mysql_maxLifetime"`
	MinIdleConnections int    `mapstructure:"mysql_minIdleConnections"`
	MaxOpenConnections int    `mapstructure:"mysql_maxOpenConnections"`
}

const currentfilepath = "internal/infrastructure/mysql/mysql.go"

func DatabaseInit(v *viper.Viper) *sql.DB {
	var mysqlConfig MysqlConf
	err := v.Unmarshal(&mysqlConfig)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, fmt.Sprintf("failed init database mysql : %s", err.Error()), nil)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DbName)
	db, err := sql.Open("mysql", dsn)
	// if there is an error opening the connection, handle it
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelFatal, "", fmt.Errorf("Cannot conenct to database : %s", err.Error()))
		panic(err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(mysqlConfig.MinIdleConnections)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(mysqlConfig.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	maxLifeTime := time.Duration(mysqlConfig.MaxLifetime) * time.Second
	db.SetConnMaxLifetime(maxLifeTime)

	if err := db.Ping(); err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "⇨ MySQL status is disconnected", err)
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "⇨ MySQL status is connected", nil)

	return db
}

func CloseDatabaseConnection(db *sql.DB) {
	db.Close()

}
