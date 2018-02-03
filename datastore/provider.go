package datastore

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type StoreResult struct {
	Data          interface{}
	ProblemDetail *models.ProblemDetail
}

//type StoreChannel chan StoreResult

type Provider interface {
	Connect() error
	Init()
	DropDatabase() error
	ApiStore
	BlockUserStore
	BotStore
	DeviceStore
	MessageStore
	RoomStore
	RoomUserStore
	MessageStore
	DeviceStore
	SettingStore
	SubscriptionStore
	UserStore
}

func GetProvider() Provider {
	var provider Provider
	switch utils.Cfg.Datastore.Provider {
	case "sqlite":
		provider = &sqliteProvider{
			sqlitePath: utils.Cfg.Datastore.SqlitePath,
		}
	case "mysql":
		provider = &mysqlProvider{
			user:              utils.Cfg.Datastore.User,
			password:          utils.Cfg.Datastore.Password,
			database:          utils.Cfg.Datastore.Database,
			masterSi:          utils.Cfg.Datastore.Master,
			replicaSis:        utils.Cfg.Datastore.Replicas,
			maxIdleConnection: utils.Cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: utils.Cfg.Datastore.MaxOpenConnection,
			trace:             false,
		}
	case "gcpSql":
		provider = &gcpSqlProvider{
			user:              utils.Cfg.Datastore.User,
			password:          utils.Cfg.Datastore.Password,
			database:          utils.Cfg.Datastore.Database,
			masterSi:          utils.Cfg.Datastore.Master,
			replicaSis:        utils.Cfg.Datastore.Replicas,
			maxIdleConnection: utils.Cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: utils.Cfg.Datastore.MaxOpenConnection,
			trace:             true,
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

func fatal(err error) {
	utils.AppLogger.Error("",
		zap.String("msg", err.Error()),
	)
	os.Exit(0)
}
