package datastore

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"

	gorp "gopkg.in/gorp.v2"

	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-sql-driver/mysql"
)

type MysqlProvider struct {
	user              string
	password          string
	database          string
	masterHost        string
	masterPort        string
	maxIdleConnection string
	maxOpenConnection string
	useSSL            string
}

func (provider MysqlProvider) Connect() error {
	if dbMap == nil {
		datasource := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			provider.user,
			provider.password,
			provider.masterHost,
			provider.masterPort,
			provider.database)
		maxIdleConnection, err := strconv.Atoi(provider.maxIdleConnection)
		if err != nil {
			log.Println(err, "RDB_MAX_IDLE_CONNECTION error")
		}
		maxOpenConnection, err := strconv.Atoi(provider.maxOpenConnection)
		if err != nil {
			log.Println(err, "RDB_MAX_OPEN_CONNECTION error")
		}

		dbMap, err = mysqlSetupConnection(
			"master",
			"mysql",
			provider.database,
			datasource,
			maxIdleConnection,
			maxOpenConnection,
			provider.useSSL,
			false)

		if err != nil {
			utils.AppLogger.Error("",
				zap.String("msg", err.Error()),
			)
			os.Exit(0)
		}
	}
	return nil
}

func (provider MysqlProvider) Init() {
	provider.CreateApiStore()
	provider.CreateUserStore()
	provider.CreateRoomStore()
	provider.CreateRoomUserStore()
	provider.CreateMessageStore()
	provider.CreateDeviceStore()
	provider.CreateSubscriptionStore()
}

func (provider MysqlProvider) DropDatabase() error {
	datasource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		provider.user,
		provider.password,
		provider.masterHost,
		provider.masterPort,
		provider.database)
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(utils.AppendStrings("DROP DATABASE ", provider.database))
	if err != nil {
		return err
	}
	return nil
}

func mysqlSetupConnection(conType, driverName, database, datasource string, maxIdle int, maxOpen int, useSSL string, trace bool) (*gorp.DbMap, error) {
	var err error
	if useSSL == "on" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(utils.Cfg.Datastore.ServerCaPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, err
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(utils.Cfg.Datastore.ClientCertPath, utils.Cfg.Datastore.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("config", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			ServerName:         utils.Cfg.Datastore.ServerName,
			InsecureSkipVerify: false,
		})
		datasource = utils.AppendStrings(datasource, "?tls=config")
	}
	db, err := sql.Open(driverName, datasource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)

	var dbmap *gorp.DbMap

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	if trace {
		dbmap.TraceOn("", log.New(os.Stdout, "sql-trace:", log.Lmicroseconds))
	}

	return dbmap, nil
}
