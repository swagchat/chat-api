package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/olebedev/config"
)

const (
	APP_NAME    = "swagchat-api"
	API_VERSION = "v0"

	DEFAULT_SERVER_PORT          = "9000"
	DEFAULT_PROFILING            = ""
	DEFAULT_SERVER_LOGGING_LEVEL = "development"
	DEFAULT_STORAGE              = "local"
	DEFAULT_DATASTORE            = "sqlite"
	DEFAULT_MESSAGING            = ""
	DEFAULT_NOTIFICATION         = ""

	DEFAULT_REALTIMESERVER_ENDPOINT = ""

	DEFAULT_LOCAL_STORAGE_BASE_URL = "http://localhost:9000/" + API_VERSION
	DEFAULT_LOCAL_STORAGE_PATH     = "data/assets"

	DEFAULT_MYSQL_MAXIDLECONNECTION = "10"
	DEFAULT_MYSQL_MAXOPENCONNECTION = "10"

	DEFAULT_SQLITE_DATABASE_PATH = "/tmp/swagchat_test.db"

	DEFAULT_MYSQL_USESSL = ""
)

type ApiServer struct {
	Port         string
	Profiling    string
	LoggingLevel string
	Storage      string
	Datastore    string
	Messaging    string
	Notification string
}

type RealtimeServer struct {
	Endpoint string
}

type LocalStorage struct {
	BaseUrl string
	Path    string
}

type GcpStorage struct {
	ProjectId           string
	Scope               string
	JwtConfigFilepath   string
	UserUploadBucket    string
	UserUploadDirectory string
	ThumbnailBucket     string
	ThumbnailDirectory  string
}

type AwsS3 struct {
	Region              string
	AccessKeyId         string
	SecretAccessKey     string
	Acl                 string
	UserUploadBucket    string
	UserUploadDirectory string
	ThumbnailBucket     string
	ThumbnailDirectory  string
}

type Sqlite struct {
	DatabasePath string
}

type Mysql struct {
	User              string
	Password          string
	Database          string
	MasterHost        string
	MasterPort        string
	MaxIdleConnection string
	MaxOpenConnection string
	UseSSL            string
	ServerName        string
	ServerCaPath      string
	ClientCertPath    string
	ClientKeyPath     string
}

type GcpSql struct {
	ProjectId         string
	SqlInstanceId     string
	User              string
	Password          string
	Database          string
	MasterHost        string
	MasterPort        string
	MaxIdleConnection string
	MaxOpenConnection string
	UseSSL            string
	ServerCaPath      string
	ClientCertPath    string
	ClientKeyPath     string
}

type GcpPubsub struct {
	ThumbnailTopic    string
	PushMessageTopic  string
	Scope             string
	JwtConfigFilepath string
}

type AwsSns struct {
	Region              string
	AccessKeyId         string
	SecretAccessKey     string
	ApplicationArn      string
	RoomTopicNamePrefix string
}

type Config struct {
	ApiServer      *ApiServer
	RealtimeServer *RealtimeServer
	// Storage
	LocalStorage *LocalStorage
	GcpStorage   *GcpStorage
	AwsS3        *AwsS3
	// Datastore
	Sqlite *Sqlite
	Mysql  *Mysql
	GcpSql *GcpSql
	// Messaging
	GcpPubsub *GcpPubsub
	// Notification
	AwsSns *AwsSns
}

var Cfg *Config = &Config{
	ApiServer:      &ApiServer{},
	RealtimeServer: &RealtimeServer{},
	// Storage
	LocalStorage: &LocalStorage{},
	GcpStorage:   &GcpStorage{},
	AwsS3:        &AwsS3{},
	// Datastore
	Mysql:  &Mysql{},
	Sqlite: &Sqlite{},
	GcpSql: &GcpSql{},
	// Messaging
	GcpPubsub: &GcpPubsub{},
	// Notification
	AwsSns: &AwsSns{},
}

func setupConfig() {
	var err error
	log.SetFlags(log.Llongfile)
	log.SetPrefix(fmt.Sprintf("[swagchat-api][%s]", API_VERSION))

	swagchatConfig := readConfig("config/swagchat.yaml")
	if swagchatConfig == nil {
		Cfg.ApiServer.Port = DEFAULT_SERVER_PORT
		Cfg.ApiServer.LoggingLevel = DEFAULT_SERVER_LOGGING_LEVEL
		Cfg.ApiServer.Datastore = DEFAULT_DATASTORE
		Cfg.ApiServer.Storage = DEFAULT_STORAGE
		Cfg.ApiServer.Messaging = DEFAULT_MESSAGING
	} else {
		if Cfg.ApiServer.Port, err = swagchatConfig.String("apiServer.port"); err != nil {
			Cfg.ApiServer.Port = DEFAULT_SERVER_PORT
		}
		if Cfg.ApiServer.Profiling, err = swagchatConfig.String("apiServer.profiling"); err != nil {
			Cfg.ApiServer.Profiling = "off"
		}
		if Cfg.ApiServer.LoggingLevel, err = swagchatConfig.String("apiServer.loggingLevel"); err != nil {
			Cfg.ApiServer.LoggingLevel = DEFAULT_SERVER_LOGGING_LEVEL
		}
		if Cfg.ApiServer.Datastore, err = swagchatConfig.String("apiServer.datastore"); err != nil {
			Cfg.ApiServer.Datastore = DEFAULT_DATASTORE
		}
		if Cfg.ApiServer.Storage, err = swagchatConfig.String("apiServer.storage"); err != nil {
			Cfg.ApiServer.Storage = DEFAULT_STORAGE
		}
		if Cfg.ApiServer.Messaging, err = swagchatConfig.String("apiServer.messaging"); err != nil {
			Cfg.ApiServer.Messaging = DEFAULT_MESSAGING
		}
		if Cfg.ApiServer.Notification, err = swagchatConfig.String("apiServer.notification"); err != nil {
			Cfg.ApiServer.Messaging = DEFAULT_NOTIFICATION
		}
		if Cfg.RealtimeServer.Endpoint, err = swagchatConfig.String("realtimeServer.endpoint"); err != nil {
			Cfg.RealtimeServer.Endpoint = DEFAULT_REALTIMESERVER_ENDPOINT
		}
	}

	if Cfg.ApiServer.Storage == "local" {
		localStorageConfig := readConfig("config/storage/local.yaml")
		if localStorageConfig == nil {
			Cfg.LocalStorage.BaseUrl = DEFAULT_LOCAL_STORAGE_BASE_URL
			Cfg.LocalStorage.Path = DEFAULT_LOCAL_STORAGE_PATH
		} else {
			if Cfg.LocalStorage.BaseUrl, err = localStorageConfig.String("localStorage.baseUrl"); err != nil {
				Cfg.LocalStorage.BaseUrl = DEFAULT_LOCAL_STORAGE_BASE_URL
			}
			if Cfg.LocalStorage.Path, err = localStorageConfig.String("localStorage.path"); err != nil {
				Cfg.LocalStorage.Path = DEFAULT_LOCAL_STORAGE_PATH
			}
		}
	}

	if Cfg.ApiServer.Storage == "gcpStorage" {
		gcpStorageConfig := readConfig("config/storage/gcpStorage.yaml")
		if gcpStorageConfig == nil {
			os.Exit(0)
		}
		if Cfg.GcpStorage.ProjectId, err = gcpStorageConfig.String("gcpStorage.projectId"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.JwtConfigFilepath, err = gcpStorageConfig.String("gcpStorage.jwtConfigFilepath"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.Scope, err = gcpStorageConfig.String("gcpStorage.scope"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.UserUploadBucket, err = gcpStorageConfig.String("gcpStorage.userUploadBucket"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.UserUploadDirectory, err = gcpStorageConfig.String("gcpStorage.userUploadDirectory"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.ThumbnailBucket, err = gcpStorageConfig.String("gcpStorage.thumbnailBucket"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpStorage.ThumbnailDirectory, err = gcpStorageConfig.String("gcpStorage.thumbnailDirectory"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
	}

	if Cfg.ApiServer.Storage == "awsS3" {
		awsS3Config := readConfig("config/storage/awsS3.yaml")
		if awsS3Config == nil {
			os.Exit(0)
		}
		if Cfg.AwsS3.Region, err = awsS3Config.String("awsS3.region"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.AccessKeyId, err = awsS3Config.String("awsS3.accessKeyId"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.SecretAccessKey, err = awsS3Config.String("awsS3.secretAccessKey"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.Acl, err = awsS3Config.String("awsS3.acl"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.UserUploadBucket, err = awsS3Config.String("awsS3.userUploadBucket"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.UserUploadDirectory, err = awsS3Config.String("awsS3.userUploadDirectory"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.ThumbnailBucket, err = awsS3Config.String("awsS3.thumbnailBucket"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsS3.ThumbnailDirectory, err = awsS3Config.String("awsS3.thumbnailDirectory"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
	}

	if Cfg.ApiServer.Datastore == "sqlite" {
		sqliteConfig := readConfig("config/datastore/sqlite.yaml")
		if sqliteConfig == nil {
			Cfg.Sqlite.DatabasePath = DEFAULT_SQLITE_DATABASE_PATH
		} else {
			if Cfg.Sqlite.DatabasePath, err = sqliteConfig.String("sqlite.databasePath"); err != nil {
				Cfg.Sqlite.DatabasePath = DEFAULT_SQLITE_DATABASE_PATH
			}
		}
	}

	if Cfg.ApiServer.Datastore == "mysql" {
		mysqlConfig := readConfig("config/datastore/mysql.yaml")
		if mysqlConfig == nil {
			os.Exit(0)
		}
		if Cfg.Mysql.User, err = mysqlConfig.String("mysql.user"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.Mysql.Password, err = mysqlConfig.String("mysql.password"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.Mysql.Database, err = mysqlConfig.String("mysql.database"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.Mysql.MasterHost, err = mysqlConfig.String("mysql.masterHost"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.Mysql.MasterPort, err = mysqlConfig.String("mysql.masterPort"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.Mysql.MaxIdleConnection, err = mysqlConfig.String("mysql.maxIdleConnection"); err != nil {
			Cfg.Mysql.MaxIdleConnection = DEFAULT_MYSQL_MAXIDLECONNECTION
		}
		if Cfg.Mysql.MaxOpenConnection, err = mysqlConfig.String("mysql.maxOpenConnection"); err != nil {
			Cfg.Mysql.MaxOpenConnection = DEFAULT_MYSQL_MAXOPENCONNECTION
		}
		if Cfg.Mysql.UseSSL, err = mysqlConfig.String("mysql.useSSL"); err != nil {
			Cfg.Mysql.UseSSL = DEFAULT_MYSQL_USESSL
		}
		if Cfg.Mysql.UseSSL == "on" {
			if Cfg.Mysql.ServerName, err = mysqlConfig.String("mysql.serverName"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if Cfg.Mysql.ServerCaPath, err = mysqlConfig.String("mysql.serverCaPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if Cfg.Mysql.ClientCertPath, err = mysqlConfig.String("mysql.clientCertPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if Cfg.Mysql.ClientKeyPath, err = mysqlConfig.String("mysql.clientKeyPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
		}
	}

	if Cfg.ApiServer.Datastore == "gcpSql" {
		gcpSqlConfig := readConfig("config/datastore/gcpSql.yaml")
		if gcpSqlConfig == nil {
			os.Exit(0)
		}
		if Cfg.GcpSql.ProjectId, err = gcpSqlConfig.String("gcpSql.projectId"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.SqlInstanceId, err = gcpSqlConfig.String("gcpSql.sqlInstanceId"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.User, err = gcpSqlConfig.String("gcpSql.user"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.Password, err = gcpSqlConfig.String("gcpSql.password"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.Database, err = gcpSqlConfig.String("gcpSql.database"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.MasterHost, err = gcpSqlConfig.String("gcpSql.masterHost"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.MasterPort, err = gcpSqlConfig.String("gcpSql.masterPort"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpSql.MaxIdleConnection, err = gcpSqlConfig.String("gcpSql.maxIdleConnection"); err != nil {
			Cfg.GcpSql.MaxIdleConnection = DEFAULT_MYSQL_MAXIDLECONNECTION
		}
		if Cfg.GcpSql.MaxOpenConnection, err = gcpSqlConfig.String("gcpSql.maxOpenConnection"); err != nil {
			Cfg.GcpSql.MaxOpenConnection = DEFAULT_MYSQL_MAXOPENCONNECTION
		}
		if Cfg.GcpSql.UseSSL, err = gcpSqlConfig.String("gcpSql.useSSL"); err != nil {
			Cfg.GcpSql.UseSSL = DEFAULT_MYSQL_USESSL
		}
		if Cfg.GcpSql.UseSSL == "on" {
			if Cfg.GcpSql.ServerCaPath, err = gcpSqlConfig.String("gcpSql.serverCaPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if Cfg.GcpSql.ClientCertPath, err = gcpSqlConfig.String("gcpSql.clientCertPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
			if Cfg.GcpSql.ClientKeyPath, err = gcpSqlConfig.String("gcpSql.clientKeyPath"); err != nil {
				log.Println(err.Error())
				os.Exit(0)
			}
		}
	}

	if Cfg.ApiServer.Messaging == "gcpPubSub" {
		gcpPubSubConfig := readConfig("config/messaging/gcpPubSub.yaml")
		if gcpPubSubConfig == nil {
			os.Exit(0)
		}
		if Cfg.GcpPubsub.ThumbnailTopic, err = gcpPubSubConfig.String("gcpPubSub.topicNameForCreateThumbnail"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpPubsub.PushMessageTopic, err = gcpPubSubConfig.String("gcpPubSub.topicNameForPushMessage"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpPubsub.Scope, err = gcpPubSubConfig.String("gcpPubSub.scope"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.GcpPubsub.JwtConfigFilepath, err = gcpPubSubConfig.String("gcpPubSub.jwtConfigFilepath"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
	}

	if Cfg.ApiServer.Notification == "awsSns" {
		awsSnsConfig := readConfig("config/notification/awsSns.yaml")
		log.Println("%#v", awsSnsConfig)
		if awsSnsConfig == nil {
			os.Exit(0)
		}
		if Cfg.AwsSns.Region, err = awsSnsConfig.String("awsSns.region"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsSns.AccessKeyId, err = awsSnsConfig.String("awsSns.accessKeyId"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsSns.SecretAccessKey, err = awsSnsConfig.String("awsSns.secretAccessKey"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsSns.ApplicationArn, err = awsSnsConfig.String("awsSns.applicationArn"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		if Cfg.AwsSns.RoomTopicNamePrefix, err = awsSnsConfig.String("awsSns.roomTopicNamePrefix"); err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
	}
}

func readConfig(filename string) *config.Config {
	absFilePath, _ := filepath.Abs(filename)
	bytes, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		return nil
	}
	yamlString := string(bytes)
	config, err := config.ParseYaml(yamlString)
	if err != nil {
		return nil
	}
	config.Env()
	return config
}

func settingFlag() {
	var stringValue string

	if stringValue = *flag.String("apiserver-port", "", "API server port.  Env: APISERVER_PORT"); stringValue != "" {
		Cfg.ApiServer.Port = stringValue
	}
	if stringValue = *flag.String("apiserver-profiling", "", "API server profiling  Env: APISERVER_PFOFILING"); stringValue != "" {
		Cfg.ApiServer.Profiling = stringValue
	}
	if stringValue = *flag.String("apiserver-loggingLevel", "", "API server logging level  Env: APISERVER_LOGGINGLEVEL"); stringValue != "" {
		Cfg.ApiServer.LoggingLevel = stringValue
	}
	if stringValue = *flag.String("apiserver-storage", "", "API server storage provider  Env: APISERVER_STORAGE"); stringValue != "" {
		Cfg.ApiServer.Storage = stringValue
	}
	if stringValue = *flag.String("apiserver-datastore", "", "API server datastore provider  Env: APISERVER_DATASTORE"); stringValue != "" {
		Cfg.ApiServer.Datastore = stringValue
	}
	if stringValue = *flag.String("apiserver-messaging", "", "API server messaging provider  Env: APISERVER_MESSAGING"); stringValue != "" {
		Cfg.ApiServer.Messaging = stringValue
	}
	if stringValue = *flag.String("apiserver-notification", "", "API server notification provider  Env: APISERVER_NOTIFICATION"); stringValue != "" {
		Cfg.ApiServer.Notification = stringValue
	}

	if stringValue = *flag.String("realtimeServer-endpoint", "", "Realtime server endpoint  Env: REALTIMESERVER_ENDPOINT"); stringValue != "" {
		Cfg.RealtimeServer.Endpoint = stringValue
	}

	if stringValue = *flag.String("localStorage-baseUrl", "", "LocalStorage base url  Env: LOCALSTORAGE_BASEURL"); stringValue != "" {
		Cfg.LocalStorage.BaseUrl = stringValue
	}
	if stringValue = *flag.String("localStorage-path", "", "LocalStorage upload path  Env: LOCALSTORAGE_PATH"); stringValue != "" {
		Cfg.LocalStorage.Path = stringValue
	}

	if stringValue = *flag.String("gcpStorage-projectId", "", "Google Cloud Storage ProjectId  Env: GCPSTORAGE_PROJECTID"); stringValue != "" {
		Cfg.GcpStorage.ProjectId = stringValue
	}
	if stringValue = *flag.String("gcpStorage-scope", "", "Google Cloud Storage Scope  Env: GCPSTORAGE_SCOPE"); stringValue != "" {
		Cfg.GcpStorage.Scope = stringValue
	}
	if stringValue = *flag.String("gcpStorage-jwtConfigFilepath", "", "Google Cloud Storage JWT config file path  Env: GCPSTORAGE_JWTCONFIGFILEPATH"); stringValue != "" {
		Cfg.GcpStorage.JwtConfigFilepath = stringValue
	}
	if stringValue = *flag.String("gcpStorage-userUploadBucket", "", "Google Cloud Storage user upload bucket name  Env: GCPSTORAGE_USERUPLOADBUCKET"); stringValue != "" {
		Cfg.GcpStorage.UserUploadBucket = stringValue
	}
	if stringValue = *flag.String("gcpStorage-userUploadDirectory", "", "Google Cloud Storage user upload directory name  Env: GCPSTORAGE_USERUPLOADDIRECTORY"); stringValue != "" {
		Cfg.GcpStorage.UserUploadDirectory = stringValue
	}
	if stringValue = *flag.String("gcpStorage-thumbnailBucket", "", "Google Cloud Storage thumbnail bucket name  Env: GCPSTORAGE_THUMBNAILBUCKET"); stringValue != "" {
		Cfg.GcpStorage.ThumbnailBucket = stringValue
	}
	if stringValue = *flag.String("gcpStorage-thumbnailDirectory", "", "Google Cloud Storage thumbnail directory name  Env: GCPSTORAGE_THUMBNAILDIRECTORY"); stringValue != "" {
		Cfg.GcpStorage.ThumbnailDirectory = stringValue
	}

	if stringValue = *flag.String("awsS3-region", "", "Amazon S3 Region  Env: AWSS3_REGION"); stringValue != "" {
		Cfg.AwsS3.Region = stringValue
	}
	if stringValue = *flag.String("awsS3-accessKeyId", "", "Amazon S3 AccessKeyId  Env: AWSS3_ACCESSKEYID"); stringValue != "" {
		Cfg.AwsS3.AccessKeyId = stringValue
	}
	if stringValue = *flag.String("awsS3-secretAccessKey", "", "Amazon S3 SecretAccessKey  Env: AWSS3_SECRETACCESSKEY"); stringValue != "" {
		Cfg.AwsS3.SecretAccessKey = stringValue
	}
	if stringValue = *flag.String("awsS3-acl", "", "Amazon S3 Acl  Env: AWSS3_ACL"); stringValue != "" {
		Cfg.AwsS3.Acl = stringValue
	}
	if stringValue = *flag.String("awsS3-userUploadBucket", "", "Amazon S3 user upload bucket name  Env: AWSS3_USERUPLOADBUCKET"); stringValue != "" {
		Cfg.AwsS3.UserUploadBucket = stringValue
	}
	if stringValue = *flag.String("awsS3-userUploadDirectory", "", "Amazon S3 user upload directory name  Env: AWSS3_USERUPLOADDIRECTORY"); stringValue != "" {
		Cfg.AwsS3.UserUploadDirectory = stringValue
	}
	if stringValue = *flag.String("awsS3-thumbnailBucket", "", "Amazon S3 thumbnail bucket name  Env: AWSS3_THUMBNAILBUCKET"); stringValue != "" {
		Cfg.AwsS3.ThumbnailBucket = stringValue
	}
	if stringValue = *flag.String("awsS3-thumbnailDirectory", "", "Amazon S3 thumbnail directory name  Env: AWSS3_THUMBNAILDIRECTORY"); stringValue != "" {
		Cfg.AwsS3.ThumbnailDirectory = stringValue
	}

	if stringValue = *flag.String("sqlite-databasePath", "", "SQLite database path  Env: SQLITE_DATABASEPATH"); stringValue != "" {
		Cfg.Sqlite.DatabasePath = stringValue
	}

	if stringValue = *flag.String("mysql-user", "", "MySQL user  Env: MYSQL_USER"); stringValue != "" {
		Cfg.Mysql.User = stringValue
	}
	if stringValue = *flag.String("mysql-password", "", "MySQL password  Env: MYSQL_PASSWORD"); stringValue != "" {
		Cfg.Mysql.Password = stringValue
	}
	if stringValue = *flag.String("mysql-database", "", "MySQL database  Env: MYSQL_DATABASE"); stringValue != "" {
		Cfg.Mysql.Database = stringValue
	}
	if stringValue = *flag.String("mysql-masterHost", "", "MySQL master host name  Env: MYSQL_MASTERHOSTNAME"); stringValue != "" {
		Cfg.Mysql.MasterHost = stringValue
	}
	if stringValue = *flag.String("mysql-masterPost", "", "MySQL master port  Env: MYSQL_MASTERPORT"); stringValue != "" {
		Cfg.Mysql.MasterPort = stringValue
	}
	if stringValue = *flag.String("mysql-maxIdleConnection", "", "MySQL max idle connection number  Env: MYSQL_MAXIDLECONNECTION"); stringValue != "" {
		Cfg.Mysql.MaxIdleConnection = stringValue
	}
	if stringValue = *flag.String("mysql-maxOpenConnection", "", "MySQL max open connection number  Env: MYSQL_MAXOPENCONNECTION"); stringValue != "" {
		Cfg.Mysql.MaxOpenConnection = stringValue
	}
	if stringValue = *flag.String("mysql-useSsl", "", "MySQL use ssl  Env: MYSQL_USESSL"); stringValue != "" {
		Cfg.Mysql.UseSSL = stringValue
	}
	if stringValue = *flag.String("mysql-servername", "", "MySQL server name  Env: MYSQL_SERVERNAME"); stringValue != "" {
		Cfg.Mysql.ServerName = stringValue
	}
	if stringValue = *flag.String("mysql-serverCaPath", "", "MySQL server CA path  Env: MYSQL_SERVERCAPATH"); stringValue != "" {
		Cfg.Mysql.ServerCaPath = stringValue
	}
	if stringValue = *flag.String("mysql-clientCertPath", "", "MySQL client cert path  Env: MYSQL_CLIENTCERTPATH"); stringValue != "" {
		Cfg.Mysql.ClientCertPath = stringValue
	}
	if stringValue = *flag.String("mysql-clientKeyPath", "", "MySQL client key path  Env: MYSQL_CLIENTKEYPATH"); stringValue != "" {
		Cfg.Mysql.ClientKeyPath = stringValue
	}

	if stringValue = *flag.String("gcpSql-projectId", "", "Google Cloud SQL projectId  Env: GCPSQL_PROJECTID"); stringValue != "" {
		Cfg.GcpSql.ProjectId = stringValue
	}
	if stringValue = *flag.String("gcpSql-sqlInstanceId", "", "Google Cloud SQL Sql InstanceId  Env: "); stringValue != "" {
		Cfg.GcpSql.SqlInstanceId = stringValue
	}
	if stringValue = *flag.String("gcpSql-user", "", "Google Cloud SQL user  Env: GCPSQL_USER"); stringValue != "" {
		Cfg.GcpSql.User = stringValue
	}
	if stringValue = *flag.String("gcpSql-password", "", "Google Cloud SQL password  Env: GCPSQL_PASSWORD"); stringValue != "" {
		Cfg.GcpSql.Password = stringValue
	}
	if stringValue = *flag.String("gcpSql-database", "", "Google Cloud SQL database  Env: GCPSQL_DATABASE"); stringValue != "" {
		Cfg.GcpSql.Database = stringValue
	}
	if stringValue = *flag.String("gcpSql-masterHost", "", "Google Cloud SQL master host name  Env: GCPSQL_MASTERHOST"); stringValue != "" {
		Cfg.GcpSql.MasterHost = stringValue
	}
	if stringValue = *flag.String("gcpSql-masterPost", "", "Google Cloud SQL master port  Env: GCPSQL_MASTERPORT"); stringValue != "" {
		Cfg.GcpSql.MasterPort = stringValue
	}
	if stringValue = *flag.String("gcpSql-maxIdleConnection", "", "Google Cloud SQL max idle connection number  Env: GCPSQL_MAXIDLECONNECTION"); stringValue != "" {
		Cfg.GcpSql.MaxIdleConnection = stringValue
	}
	if stringValue = *flag.String("gcpSql-maxOpenConnection", "", "Google Cloud SQL max open connection number  Env: GCPSQL_MAXOPENCONNECTION"); stringValue != "" {
		Cfg.GcpSql.MaxOpenConnection = stringValue
	}
	if stringValue = *flag.String("gcpSql-useSsl", "", "Google Cloud SQL use ssl  Env: GCPSQL_USESSL"); stringValue != "" {
		Cfg.GcpSql.UseSSL = stringValue
	}
	if stringValue = *flag.String("gcpSql-serverCaPath", "", "Google Cloud SQL server CA path  Env: GCPSQL_SERVERCAPATH"); stringValue != "" {
		Cfg.GcpSql.ServerCaPath = stringValue
	}
	if stringValue = *flag.String("gcpSql-clientCertPath", "", "Google Cloud SQL client cert path  Env: GCPSQL_CLIENTCERTPATH"); stringValue != "" {
		Cfg.GcpSql.ClientCertPath = stringValue
	}
	if stringValue = *flag.String("gcpSql-clientKeyPath", "", "Google Cloud SQL client key path  Env: GCPSQL_CLIENTKEYPATH"); stringValue != "" {
		Cfg.GcpSql.ClientKeyPath = stringValue
	}

	if stringValue = *flag.String("gcpPubsub-thumbnailTopic", "", "Google Cloud Pub/Sub thumbnail topic name  Env: GCPPUBSUB_THUMBNAILTOPIC"); stringValue != "" {
		Cfg.GcpPubsub.ThumbnailTopic = stringValue
	}
	if stringValue = *flag.String("gcpPubsub-pushMessageTopic", "", "Google Cloud Pub/Sub push message topic name  Env: GCPPUBSUB_PUSHMESSAGETOPIC"); stringValue != "" {
		Cfg.GcpPubsub.PushMessageTopic = stringValue
	}
	if stringValue = *flag.String("gcpPubsub-scope", "", "Google Cloud Pub/Sub scope  Env: GCPPUBSUB_SCOPE"); stringValue != "" {
		Cfg.GcpPubsub.Scope = stringValue
	}
	if stringValue = *flag.String("gcpPubsub-jwtConfigFilepath", "", "Google Cloud Pub/Sub JWT config file path  Env: GCPPUBSUB_JWTCONFIGFILEPATH"); stringValue != "" {
		Cfg.GcpPubsub.JwtConfigFilepath = stringValue
	}

	if stringValue = *flag.String("awssns-region", "", "Amazon SNS Region  Env: AWSSNS_REGION"); stringValue != "" {
		Cfg.AwsSns.Region = stringValue
	}
	if stringValue = *flag.String("awssns-accessKeyId", "", "Amazon SNS AccessKeyId  Env: AWSSNS_ACCESSKEYID"); stringValue != "" {
		Cfg.AwsSns.AccessKeyId = stringValue
	}
	if stringValue = *flag.String("awssns-secretAccessKey", "", "Amazon SNS SecretAccessKey  Env: AWSSNS_SECRETACCESSKEY"); stringValue != "" {
		Cfg.AwsSns.SecretAccessKey = stringValue
	}
	if stringValue = *flag.String("awssns-applicationArn", "", "Amazon SNS ApplicationArn  Env: AWSSNS_APPLICATIONARN"); stringValue != "" {
		Cfg.AwsSns.ApplicationArn = stringValue
	}
	if stringValue = *flag.String("awssns-roomTopicNamePrefix", "", "Amazon SNS RoomTopicNamePrefix  Env: AWSSNS_ROOMTOPICNAMEPREFIX"); stringValue != "" {
		Cfg.AwsSns.RoomTopicNamePrefix = stringValue
	}

	flag.Parse()
}
