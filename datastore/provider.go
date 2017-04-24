package datastore

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/pkg/errors"
)

type StoreResult struct {
	Data          interface{}
	ProblemDetail *models.ProblemDetail
}

type StoreChannel chan StoreResult

type Provider interface {
	Connect() error
	Init()
	DropDatabase() error
	UserStore
	RoomStore
	RoomUserStore
	MessageStore
	DeviceStore
	SubscriptionStore
}

func GetProvider() Provider {
	var provider Provider
	switch utils.Cfg.ApiServer.Datastore {
	case "sqlite":
		provider = &SqliteProvider{
			databasePath: utils.Cfg.Sqlite.DatabasePath,
		}
	case "mysql":
		provider = &MysqlProvider{
			user:              utils.Cfg.Mysql.User,
			password:          utils.Cfg.Mysql.Password,
			database:          utils.Cfg.Mysql.Database,
			masterHost:        utils.Cfg.Mysql.MasterHost,
			masterPort:        utils.Cfg.Mysql.MasterPort,
			maxIdleConnection: utils.Cfg.Mysql.MaxIdleConnection,
			maxOpenConnection: utils.Cfg.Mysql.MaxOpenConnection,
			useSSL:            utils.Cfg.Mysql.UseSSL,
		}
	case "gcpSql":
		provider = &GcpSqlProvider{
			user:              utils.Cfg.GcpSql.User,
			password:          utils.Cfg.GcpSql.Password,
			database:          utils.Cfg.GcpSql.Database,
			masterHost:        utils.Cfg.GcpSql.MasterHost,
			masterPort:        utils.Cfg.GcpSql.MasterPort,
			maxIdleConnection: utils.Cfg.GcpSql.MaxIdleConnection,
			maxOpenConnection: utils.Cfg.GcpSql.MaxOpenConnection,
			useSSL:            utils.Cfg.GcpSql.UseSSL,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.ApiServer.Datastore is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}

func createProblemDetail(title string, err error) *models.ProblemDetail {
	if err == nil {
		err = errors.New("")
	}
	return &models.ProblemDetail{
		Title:     title,
		Status:    http.StatusInternalServerError,
		ErrorName: models.ERROR_NAME_DATABASE_ERROR,
		Detail:    err.Error(),
		Error:     errors.Wrap(err, title),
	}
}
