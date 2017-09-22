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

	gorp "gopkg.in/gorp.v2"

	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-sql-driver/mysql"
)

type gcpSqlProvider struct {
	user              string
	password          string
	database          string
	masterHost        string
	masterPort        string
	slaveHost         string
	slavePort         string
	maxIdleConnection string
	maxOpenConnection string
	useSSL            bool
	trace             bool
}

func (p *gcpSqlProvider) Connect() error {
	rs := RdbStoreInstance()
	if rs.Master() == nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.masterHost,
			p.masterPort,
			p.database)
		db, err := p.openDb(ds, p.useSSL)
		if err != nil {
			fatal(err)
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
		if p.trace {
			master.TraceOn("", log.New(os.Stdout, "sql-trace:", log.Lmicroseconds))
		}

		rs.SetMaster(master)
	}
	if p.slaveHost != "" && p.slavePort != "" && rs.Slave() == nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.slaveHost,
			p.slavePort,
			p.database)
		db, err := p.openDb(ds, p.useSSL)
		if err != nil {
			fatal(err)
		}

		mic, err := strconv.Atoi(p.maxIdleConnection)
		if err == nil {
			db.SetMaxIdleConns(mic)
		}
		moc, err := strconv.Atoi(p.maxOpenConnection)
		if err == nil {
			db.SetMaxOpenConns(moc)
		}

		db.SetMaxIdleConns(mic)
		db.SetMaxOpenConns(moc)

		var slave *gorp.DbMap
		slave = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
		if p.trace {
			slave.TraceOn("", log.New(os.Stdout, "sql-trace:", log.Lmicroseconds))
		}

		rs.SetSlave(slave)
	}
	return nil
}

func (p *gcpSqlProvider) Init() {
	p.CreateApiStore()
	p.CreateUserStore()
	p.CreateBlockUserStore()
	p.CreateRoomStore()
	p.CreateRoomUserStore()
	p.CreateMessageStore()
	p.CreateDeviceStore()
	p.CreateSubscriptionStore()
}

func (p *gcpSqlProvider) DropDatabase() error {
	rs := RdbStoreInstance()
	if rs.Master() != nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.masterHost,
			p.masterPort,
			p.database)
		db, err := p.openDb(ds, p.useSSL)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec(utils.AppendStrings("DROP DATABASE ", p.database))
		if err != nil {
			return err
		}
	}
	if p.slaveHost != "" && p.slavePort != "" && rs.Slave() == nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.slaveHost,
			p.slavePort,
			p.database)
		db, err := p.openDb(ds, p.useSSL)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec(utils.AppendStrings("DROP DATABASE ", p.database))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *gcpSqlProvider) openDb(dataSource string, useSSL bool) (*sql.DB, error) {
	var err error
	if useSSL {
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
			ServerName:         utils.AppendStrings(utils.Cfg.Datastore.GcpProjectId, ":", utils.Cfg.Datastore.ServerName),
			InsecureSkipVerify: false,
		})
		dataSource = utils.AppendStrings(dataSource, "?tls=config")
	}
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
