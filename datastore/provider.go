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
	AssetStore
	BlockUserStore
	BotStore
	DeviceStore
	MessageStore
	RoomStore
	RoomUserStore
	SettingStore
	SubscriptionStore
	UserStore
}

func GetProvider() Provider {
	cfg := utils.GetConfig()

	var provider Provider
	switch cfg.Datastore.Provider {
	case "sqlite":
		provider = &sqliteProvider{
			sqlitePath: cfg.Datastore.SQLite.Path,
			trace:      true,
		}
	case "mysql":
		provider = &mysqlProvider{
			user:              cfg.Datastore.User,
			password:          cfg.Datastore.Password,
			database:          cfg.Datastore.Database,
			masterSi:          cfg.Datastore.Master,
			replicaSis:        cfg.Datastore.Replicas,
			maxIdleConnection: cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: cfg.Datastore.MaxOpenConnection,
			trace:             false,
		}
	case "gcpSql":
		provider = &gcpSqlProvider{
			user:              cfg.Datastore.User,
			password:          cfg.Datastore.Password,
			database:          cfg.Datastore.Database,
			masterSi:          cfg.Datastore.Master,
			replicaSis:        cfg.Datastore.Replicas,
			maxIdleConnection: cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: cfg.Datastore.MaxOpenConnection,
			trace:             false,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "cfg.ApiServer.Datastore is incorrect"),
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
