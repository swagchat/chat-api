package datastore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"

	gorp "gopkg.in/gorp.v2"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
)

type mysqlProvider struct {
	ctx               context.Context
	user              string
	password          string
	database          string
	masterSi          *config.ServerInfo
	replicaSis        []*config.ServerInfo
	maxIdleConnection int
	maxOpenConnection int
	enableLogging     bool
}

func (p *mysqlProvider) Connect(dsCfg *config.Datastore) error {
	if _, ok := rdbStores[dsCfg.Database]; ok {
		return nil
	}

	rs := &rdbStore{}
	if rs.master() == nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.masterSi.Host,
			p.masterSi.Port,
			p.database)
		db, err := p.openDb(ds, p.masterSi)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("Failed to connect database. %s %s", dsCfg.Provider, ds))
			logger.Error(err.Error())
			return err
		}
		logger.Info(fmt.Sprintf("Connected database. %s %s", dsCfg.Provider, ds))

		db.SetMaxIdleConns(p.maxIdleConnection)
		db.SetMaxOpenConns(p.maxOpenConnection)

		var master *gorp.DbMap
		master = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
		if p.enableLogging {
			master.TraceOn("[master]", logger.Logger())
		}
		rs.setMaster(master)
	}

	for _, replicaSi := range p.replicaSis {
		if replicaSi.Host != "" && replicaSi.Port != "" && rs.replica() == nil {
			ds := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s",
				p.user,
				p.password,
				replicaSi.Host,
				replicaSi.Port,
				p.database)
			db, err := p.openDb(ds, replicaSi)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Failed to connect database. %s %s", dsCfg.Provider, ds))
				logger.Error(err.Error())
				return err
			}
			logger.Info(fmt.Sprintf("Connected database. %s %s", dsCfg.Provider, ds))

			db.SetMaxIdleConns(p.maxIdleConnection)
			db.SetMaxOpenConns(p.maxOpenConnection)

			var replica *gorp.DbMap
			replica = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
			if p.enableLogging {
				replica.TraceOn("[replica]", logger.Logger())
			}
			rs.setReplica(replica)
		}
	}
	rdbStores[dsCfg.Database] = rs
	return nil
}

func (p *mysqlProvider) CreateTables() {
	p.createAppClientStore()
	p.createAssetStore()
	p.createBlockUserStore()
	p.createDeviceStore()
	p.createMessageStore()
	p.createRoomStore()
	p.createRoomUserStore()
	p.createSettingStore()
	p.createSubscriptionStore()
	p.createUserStore()
	p.createUserRoleStore()
	p.createWebhookStore()
}

func (p *mysqlProvider) DropDatabase() error {
	master := RdbStore(p.database).master()
	if master != nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.masterSi.Host,
			p.masterSi.Port,
			p.database)
		db, err := p.openDb(ds, p.masterSi)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		defer db.Close()
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s", p.database))
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return nil
}

func (p *mysqlProvider) openDb(dataSource string, si *config.ServerInfo) (*sql.DB, error) {
	var err error
	if si.ServerName != "" && si.ServerCaPath != "" && si.ClientCertPath != "" && si.ClientKeyPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(si.ServerCaPath)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			err = errors.New("Failed to append certs from PEM")
			logger.Error(err.Error())
			return nil, err
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(si.ClientCertPath, si.ClientKeyPath)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("config", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			ServerName:         si.ServerName,
			InsecureSkipVerify: false,
		})
		dataSource = fmt.Sprintf("%s?tls=config", dataSource)
	}
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return db, nil
}

func (p *mysqlProvider) Close() {
	for database, rdbStore := range rdbStores {
		if rdbStore != nil {
			master := rdbStore.masterDbMap
			if master != nil {
				close(database, master.Db)
			}
			for _, replica := range rdbStores[p.database].replicaDbMaps {
				close(database, replica.Db)
			}
			delete(rdbStores, database)
		}
	}
}
