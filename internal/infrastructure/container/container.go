package container

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/mysql"
	redisclient "github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/redis"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/viper"
)

var v *viper.Viper

const currentfilepath = "internal/infrastructure/container/container.go"

type (
	Container struct {
		Mysqldb   *sql.DB
		Apps      *Apps
		Logger    *Logger
		Redis     *redis.Client
		ReqResAPI *ReqResAPI
	}

	Logger struct {
		Log     zerolog.Logger
		Path    string `mapstructure:"log_path"`
		LogFile string `mapstructure:"log_file"`
	}

	Apps struct {
		Name      string `mapstructure:"appName"`
		Host      string `mapstructure:"host"`
		Version   string `mapstructure:"version"`
		Address   string `mapstructure:"address"`
		HttpPort  int    `mapstructure:"httpport"`
		SecretJwt string `mapstructure:"secretJwt"`
	}

	ReqResAPI struct {
		URL       string `mapstructure:"reqres_uri"`
		TimeOut   int    `mapstructure:"reqres_timeout"`
		Debugging bool   `mapstructure:"reqres_debugging"`
	}
)

func loadEnv() {
	projectDirName := "go-simple-user-crud"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	v.SetConfigFile(string(rootPath) + `/.env`)
}

func init() {
	v = viper.New()

	v.AutomaticEnv()
	loadEnv()

	path, err := os.Executable()
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("os.Executable panic : %s", err.Error()))
	}

	dir := filepath.Dir(path)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("failed read config : %s", err.Error()))
	}

	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read configuration file", nil)
}

func AppsInit(v *viper.Viper) (apps Apps) {
	err := v.Unmarshal(&apps)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when unmarshal configuration file : ", err.Error()))
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed when unmarshal configuration file", nil)
	return
}

func LoggerInit(v *viper.Viper) (logger Logger) {
	err := v.Unmarshal(&logger)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when unmarshal configuration file : ", err.Error()))
	}

	var stdout io.Writer = os.Stdout
	// var multi zerolog.LevelWriter = os.Stdout
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if logger.LogFile == "ON" {
		path := fmt.Sprintf("%s/logs/request.log", helper.ProjectRootPath)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664)
		if err != nil {
			log.Error().Err(err)
		}
		// Create a multi writer with both the console and file writers
		stdout = zerolog.MultiLevelWriter(os.Stdout, file)

	}

	return Logger{
		Log: zerolog.New(stdout).With().Caller().Timestamp().Logger(),
	}
}

func ReqResAPIInit(v *viper.Viper) (reqresapi ReqResAPI) {
	err := v.Unmarshal(&reqresapi)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when unmarshal configuration file : ", err.Error()))
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed when unmarshal configuration file", nil)
	return
}

// containters = apps,mysql,logger,redis
func InitContainer(containters ...string) *Container {
	newStrContainer := strings.Join(containters, ",")
	var cont Container
	errGroup, _ := errgroup.WithContext(context.Background())

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "apps") || len(containters) == 0 {
			apps := AppsInit(v)
			cont.Apps = &apps
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "mysql") || len(containters) == 0 {
			mysqldb := mysql.DatabaseInit(v)
			cont.Mysqldb = mysqldb
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "log") || len(containters) == 0 {
			logger := LoggerInit(v)
			cont.Logger = &logger
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "redis") || len(containters) == 0 {
			redisClient := redisclient.NewRedisClient(v)
			cont.Redis = redisClient
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "reqresapi") || len(containters) == 0 {
			reqresapi := ReqResAPIInit(v)
			cont.ReqResAPI = &reqresapi
			return
		}
		return nil
	})

	_ = errGroup.Wait()

	return &cont
}

func GetEnv() *viper.Viper {
	return v
}
