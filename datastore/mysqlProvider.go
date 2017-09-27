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

	"github.com/go-sql-driver/mysql"
	"github.com/swagchat/chat-api/utils"
)

type mysqlProvider struct {
	user              string
	password          string
	database          string
	masterSi          *utils.ServerInfo
	replicaSis        []*utils.ServerInfo
	maxIdleConnection string
	maxOpenConnection string

	trace bool
}

func (p *mysqlProvider) Connect() error {
	rs := RdbStoreInstance()
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

			rs.setReplica(slave)
		}
	}
	return nil
}

func (p *mysqlProvider) Init() {
	p.CreateApiStore()
	p.CreateUserStore()
	p.CreateBlockUserStore()
	p.CreateRoomStore()
	p.CreateRoomUserStore()
	p.CreateMessageStore()
	p.CreateDeviceStore()
	p.CreateSubscriptionStore()
}

func (p *mysqlProvider) DropDatabase() error {
	rs := RdbStoreInstance()
	if rs.master() != nil {
		ds := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			p.user,
			p.password,
			p.masterSi.Host,
			p.masterSi.Port,
			p.database)
		db, err := p.openDb(ds, p.masterSi)
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

func (p *mysqlProvider) openDb(dataSource string, si *utils.ServerInfo) (*sql.DB, error) {
	var err error
	if si.ServerName != "" && si.ServerCaPath != "" && si.ClientCertPath != "" && si.ClientKeyPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(si.ServerCaPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, err
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(si.ClientCertPath, si.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("config", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			ServerName:         si.ServerName,
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
