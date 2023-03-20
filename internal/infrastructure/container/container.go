package container

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"

	mysqlclient "github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/mysql"
	redisclient "github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/redis"
	"github.com/go-redis/redis/v8"
	"github.com/mashingan/smapping"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/sync/errgroup"
)

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
		Path    string `env:"log_path"`
		LogFile string `env:"log_file"`
	}

	Apps struct {
		Name      string `env:"apps_appName"`
		Host      string `env:"apps_host"`
		Version   string `env:"apps_version"`
		Address   string `env:"apps_address"`
		HttpPort  int    `env:"apps_httpport"`
		SecretJwt string `env:"apps_secretJwt"`
	}

	ReqResAPI struct {
		URL       string `env:"reqres_uri"`
		TimeOut   int    `env:"reqres_timeout"`
		Debugging bool   `env:"reqres_debugging"`
	}
)

var mapsEnv = smapping.Mapped{}

func init() {
	err := godotenv.Load(fmt.Sprintf("%s/.env", helper.ProjectRootPath))
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when loadenv : ", err.Error()))
	}

	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "apps_") || strings.HasPrefix(v, "mysql_") || strings.HasPrefix(v, "redis_") || strings.HasPrefix(v, "log_") || strings.HasPrefix(v, "reqres_") {
			strs := strings.Split(v, "=")
			mapsEnv[strs[0]] = strs[1]
		}
	}

	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read environment variable", nil)
}

func AppsInit() Apps {
	var appsConf = Apps{
		Name:      utils.EnvString("apps_appName"),
		Host:      utils.EnvString("apps_host"),
		Version:   utils.EnvString("apps_version"),
		Address:   utils.EnvString("apps_address"),
		HttpPort:  utils.EnvInt("apps_httpport"),
		SecretJwt: utils.EnvString("apps_secretJwt"),
	}

	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read appsConf", nil)
	return appsConf
}

func LoggerInit() Logger {
	var loggerConf = Logger{
		Path:    utils.EnvString("log_path"),
		LogFile: utils.EnvString("log_file"),
	}

	err := smapping.FillStructByTags(&loggerConf, mapsEnv, "env")
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when read loggerConf : ", err.Error()))
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed when read loggerConf", nil)

	var stdout io.Writer = os.Stdout
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if loggerConf.LogFile == "ON" {
		path := fmt.Sprintf("%s%s", helper.ProjectRootPath, loggerConf.Path)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664)
		if err != nil {
			helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when setting loggerConf : ", err.Error()))
		}
		// Create a multi writer with both the console and file writers
		stdout = zerolog.MultiLevelWriter(os.Stdout, file)

	}

	loggerConf.Log = zerolog.New(stdout).With().Caller().Timestamp().Logger()
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read loggerConf", nil)
	return loggerConf
}

func ReqResAPIInit() ReqResAPI {
	var reqresConf = ReqResAPI{
		URL:       utils.EnvString("reqres_uri"),
		TimeOut:   utils.EnvInt("reqres_timeout"),
		Debugging: utils.EnvBool("reqres_debugging"),
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read reqresConf", nil)
	return reqresConf
}

// containters = apps,mysql,logger,redis
func InitContainer(containters ...string) *Container {
	newStrContainer := strings.Join(containters, ",")
	var cont Container
	errGroup, _ := errgroup.WithContext(context.Background())

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "apps") || len(containters) == 0 {
			apps := AppsInit()
			cont.Apps = &apps
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "mysql") || len(containters) == 0 {
			mysqldb := mysqlclient.DatabaseInit()
			cont.Mysqldb = mysqldb
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "log") || len(containters) == 0 {
			logger := LoggerInit()
			cont.Logger = &logger
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "redis") || len(containters) == 0 {
			redisClient := redisclient.NewRedisClient()
			cont.Redis = redisClient
			return
		}
		return nil
	})

	errGroup.Go(func() (err error) {
		if strings.Contains(newStrContainer, "reqresapi") || len(containters) == 0 {
			reqresapi := ReqResAPIInit()
			cont.ReqResAPI = &reqresapi
			return
		}
		return nil
	})

	_ = errGroup.Wait()

	return &cont
}
