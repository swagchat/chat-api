package datastore

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strconv"

	gorp "gopkg.in/gorp.v2"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
)

type gcpSQLProvider struct {
	user              string
	password          string
	database          string
	masterSi          *utils.ServerInfo
	replicaSis        []*utils.ServerInfo
	maxIdleConnection string
	maxOpenConnection string
	enableLogging     bool
}

func (p *gcpSQLProvider) Connect(dsCfg *utils.Datastore) error {
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
			err = errors.Wrap(err, "Database open error")
			logger.Error(err.Error())
			return err
		}

		mic, err := strconv.Atoi(p.maxIdleConnection)
		if err == nil {
			db.SetMaxIdleConns(mic)
		}

		moc, err := strconv.Atoi(p.maxOpenConnection)
		if err == nil {
			db.SetMaxOpenConns(moc)
		}

		var master *gorp.DbMap
		master = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
		if p.enableLogging {
			master.TraceOn("[master]", logger.Logger())

		}
		rs.setMaster(master)
	}

	for _, replicaSi := range p.replicaSis {
		if replicaSi.Host != "" && replicaSi.Port != "" && rs.replicaDbMaps == nil {
			ds := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s",
				p.user,
				p.password,
				replicaSi.Host,
				replicaSi.Port,
				p.database)
			db, err := p.openDb(ds, replicaSi)
			if err != nil {
				err = errors.Wrap(err, "Database open error")
				logger.Error(err.Error())
				return err
			}

			mic, err := strconv.Atoi(p.maxIdleConnection)
			if err == nil {
				db.SetMaxIdleConns(mic)
			}

			moc, err := strconv.Atoi(p.maxOpenConnection)
			if err == nil {
				db.SetMaxOpenConns(moc)
			}

			var replica *gorp.DbMap
			replica = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
			if p.enableLogging {
				replica.TraceOn("[replica]", logger.Logger())
			}
			rs.setReplica(replica)
		}
	}
	rdbStores[dsCfg.Database] = rs
	p.init()
	return nil
}

func (p *gcpSQLProvider) init() {
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

func (p *gcpSQLProvider) DropDatabase() error {
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

func (p *gcpSQLProvider) openDb(dataSource string, si *utils.ServerInfo) (*sql.DB, error) {
	var err error
	if si.ServerName != "" && si.ServerCaPath != "" && si.ClientCertPath != "" && si.ClientKeyPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(si.ServerCaPath)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			err = errors.New("Failed AppendCertsFromPEM.")
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
