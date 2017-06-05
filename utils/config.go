package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	APP_NAME      = "swagchat-api"
	API_VERSION   = "v0"
	BUILD_VERSION = "v0.6.0"

	KEY_LENGTH        = 32
	TOKEN_LENGTH      = 32
	HEADER_API_KEY    = "X-SwagChat-Api-Key"
	HEADER_API_SECRET = "X-SwagChat-Api-Secret"
	HEADER_USER_ID    = "X-SwagChat-User-Id"
)

var (
	Cfg           *Config
	IsShowVersion bool
)

type Config struct {
	Version        string
	Port           string
	Profiling      bool
	ErrorLogging   bool `yaml:"errorLogging"`
	Logging        *Logging
	Storage        *Storage
	Datastore      *Datastore
	Messaging      *Messaging
	Notification   *Notification
	RealtimeServer *RealtimeServer `yaml:"realtimeServer"`
}

type Logging struct {
	Level string
}

type Storage struct {
	Provider           string
	UploadBucket       string `yaml:"uploadBucket"`
	UploadDirectory    string `yaml:"uploadDirectory"`
	ThumbnailBucket    string `yaml:"thumbnailBucket"`
	ThumbnailDirectory string `yaml:"thumbnailDirectory"`

	// Local
	BaseUrl   string `yaml:"baseUrl"`
	LocalPath string `yaml:"localPath"`

	// GCP Storage
	GcpProjectId string `yaml:"gcpProjectId"`
	GcpJwtPath   string `yaml:"gcpJwtPath"`

	// AWS S3
	AwsRegion          string `yaml:"awsRegion"`
	AwsAccessKeyId     string `yaml:"awsAccessKeyId"`
	AwsSecretAccessKey string `yaml:"AwsSecretAccessKey"`
}

type Datastore struct {
	Provider        string
	TableNamePrefix string `yaml:"tableNamePrefix"`

	// SQLite
	SqlitePath string `yaml:"sqlitePath"`

	// MySQL, GCP SQL
	User              string
	Password          string
	Database          string
	MasterHost        string `yaml:"masterHost"`
	MasterPort        string `yaml:"masterPort"`
	MaxIdleConnection string `yaml:"maxIdleConnection"`
	MaxOpenConnection string `yaml:"maxOpenConnection"`
	UseSSL            string `yaml:"useSSL"`     // "on" or "off"
	ServerName        string `yaml:"serverName"` // For GcpSQL, set SqlInstanceId.
	ServerCaPath      string `yaml:"serverCaPath"`
	ClientCertPath    string `yaml:"clientCertPath"`
	ClientKeyPath     string `yaml:"clientKeyPath"`

	// GCP SQL
	GcpProjectId string `yaml:"gcpProjectId"`
}

type Messaging struct {
	Provider       string
	ThumbnailTopic string `yaml:"thumbnailTopic"`

	// GCP Pubsub
	GcpProjectId string `yaml:"gcpProjectId"`
	GcpJwtPath   string `yaml:"gcpJwtPath"`
}

type Notification struct {
	Provider            string
	RoomTopicNamePrefix string `yaml:"roomTopicNamePrefix"`

	// AWS SNS
	AwsRegion                string `yaml:"awsRegion"`
	AwsAccessKeyId           string `yaml:"awsAccessKeyId"`
	AwsSecretAccessKey       string `yaml:"awsSecretAccessKey"`
	AwsApplicationArnIos     string `yaml:"awsApplicationArnIos"`
	AwsApplicationArnAndroid string `yaml:"awsApplicationArnAndroid"`
}

type RealtimeServer struct {
	Endpoint string
}

func setupConfig() {
	log.SetFlags(log.Llongfile)
	loadDefaultSettings()
	loadYaml()
	loadEnvironment()
	parseFlag()
}

func loadDefaultSettings() {
	logging := &Logging{
		Level: "development",
	}

	storage := &Storage{
		Provider:  "local",
		BaseUrl:   "http://localhost:9000/v0/assets",
		LocalPath: "data/assets",
	}

	datastore := &Datastore{
		Provider:   "sqlite",
		SqlitePath: "/tmp/swagchat.db",
		UseSSL:     "off",
	}

	messaging := &Messaging{}

	notification := &Notification{}

	realtimeServer := &RealtimeServer{
		Endpoint: "",
	}

	Cfg = &Config{
		Port:           "9000",
		Profiling:      false,
		ErrorLogging:   false,
		Logging:        logging,
		Storage:        storage,
		Datastore:      datastore,
		Messaging:      messaging,
		Notification:   notification,
		RealtimeServer: realtimeServer,
	}
}

func loadYaml() {
	buf, _ := ioutil.ReadFile("config/swagchat.yaml")
	yaml.Unmarshal(buf, Cfg)
}

func loadEnvironment() {
	var v string

	if v = os.Getenv("SC_PORT"); v != "" {
		Cfg.Port = v
	}
	if v = os.Getenv("SC_PROFILING"); v != "" {
		if v == "on" {
			Cfg.Profiling = true
		} else {
			Cfg.Profiling = false
		}
	}
	if v = os.Getenv("SC_ERROR_LOGGING"); v != "" {
		if v == "on" {
			Cfg.ErrorLogging = true
		} else {
			Cfg.ErrorLogging = false
		}
	}

	// Logging
	if v = os.Getenv("SC_LOGGING_LEVEL"); v != "" {
		Cfg.Logging.Level = v
	}

	// Storage
	if v = os.Getenv("SC_STORAGE_PROVIDER"); v != "" {
		Cfg.Storage.Provider = v
	}
	if v = os.Getenv("SC_STORAGE_UPLOAD_BUCKET"); v != "" {
		Cfg.Storage.UploadBucket = v
	}
	if v = os.Getenv("SC_STORAGE_UPLOAD_DIRECTORY"); v != "" {
		Cfg.Storage.UploadDirectory = v
	}
	if v = os.Getenv("SC_STORAGE_THUMBNAIL_BUCKET"); v != "" {
		Cfg.Storage.ThumbnailBucket = v
	}
	if v = os.Getenv("SC_STORAGE_THUMBNAIL_DIRECTORY"); v != "" {
		Cfg.Storage.ThumbnailDirectory = v
	}

	// Storage - Local
	if v = os.Getenv("SC_STORAGE_BASE_URL"); v != "" {
		Cfg.Storage.BaseUrl = v
	}
	if v = os.Getenv("SC_STORAGE_LOCAL_PATH"); v != "" {
		Cfg.Storage.LocalPath = v
	}

	// Storage - GCP Storage
	if v = os.Getenv("SC_STORAGE_GCP_PROJECT_ID"); v != "" {
		Cfg.Storage.GcpProjectId = v
	}
	if v = os.Getenv("SC_STORAGE_GCP_JWT_PATH"); v != "" {
		Cfg.Storage.GcpJwtPath = v
	}

	// Storage - AWS S3
	if v = os.Getenv("SC_STORAGE_AWS_REGION"); v != "" {
		Cfg.Storage.AwsRegion = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_ACCESS_KEY_ID"); v != "" {
		Cfg.Storage.AwsAccessKeyId = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_SECRET_ACCESS_KEY"); v != "" {
		Cfg.Storage.AwsSecretAccessKey = v
	}

	// Datastore
	if v = os.Getenv("SC_DATASTORE_PROVIDER"); v != "" {
		Cfg.Datastore.Provider = v
	}
	if v = os.Getenv("SC_DATASTORE_TABLE_NAME_PREFIX"); v != "" {
		Cfg.Datastore.TableNamePrefix = v
	}

	// Datastore - SQLite
	if v = os.Getenv("SC_DATASTORE_SQLITE_PATH"); v != "" {
		Cfg.Datastore.SqlitePath = v
	}

	// Datastore - MySQL, GCP SQL
	if v = os.Getenv("SC_DATASTORE_USER"); v != "" {
		Cfg.Datastore.User = v
	}
	if v = os.Getenv("SC_DATASTORE_PASSWORD"); v != "" {
		Cfg.Datastore.Password = v
	}
	if v = os.Getenv("SC_DATASTORE_DATABASE"); v != "" {
		Cfg.Datastore.Database = v
	}
	if v = os.Getenv("SC_DATASTORE_MASTER_HOST"); v != "" {
		Cfg.Datastore.MasterHost = v
	}
	if v = os.Getenv("SC_DATASTORE_MASTER_PORT"); v != "" {
		Cfg.Datastore.MasterPort = v
	}
	if v = os.Getenv("SC_DATASTORE_MAX_IDLE_CONNECTION"); v != "" {
		Cfg.Datastore.MaxIdleConnection = v
	}
	if v = os.Getenv("SC_DATASTORE_MAX_OPEN_CONNECTION"); v != "" {
		Cfg.Datastore.MaxOpenConnection = v
	}
	if v = os.Getenv("SC_DATASTORE_USE_SSL"); v != "" {
		Cfg.Datastore.UseSSL = v
	}
	if v = os.Getenv("SC_DATASTORE_SERVER_NAME"); v != "" {
		Cfg.Datastore.ServerName = v
	}
	if v = os.Getenv("SC_DATASTORE_SERVER_CA_PATH"); v != "" {
		Cfg.Datastore.ServerCaPath = v
	}
	if v = os.Getenv("SC_DATASTORE_CLIENT_CERT_PATH"); v != "" {
		Cfg.Datastore.ClientCertPath = v
	}
	if v = os.Getenv("SC_DATASTORE_CLIENT_KEY_PATH"); v != "" {
		Cfg.Datastore.ClientKeyPath = v
	}

	// Datastore - GCP SQL
	if v = os.Getenv("SC_DATASTORE_GCP_PROJECT_ID"); v != "" {
		Cfg.Datastore.GcpProjectId = v
	}

	// Messaging
	if v = os.Getenv("SC_MESSAGING_PROVIDER"); v != "" {
		Cfg.Messaging.Provider = v
	}
	if v = os.Getenv("SC_NOTIFICATION_ROOM_TOPIC_NAME_PREFIX"); v != "" {
		Cfg.Notification.RoomTopicNamePrefix = v
	}

	// Messaging - GCP Pubsub
	if v = os.Getenv("SC_MESSAGING_GCP_PROJECT_ID"); v != "" {
		Cfg.Messaging.GcpProjectId = v
	}
	if v = os.Getenv("SC_MESSAGING_GCP_JWT_PATH"); v != "" {
		Cfg.Messaging.GcpJwtPath = v
	}
	if v = os.Getenv("SC_MESSAGING_THUMBNAIL_TOPIC"); v != "" {
		Cfg.Messaging.ThumbnailTopic = v
	}

	// Notification
	if v = os.Getenv("SC_NOTIFICATION_PROVIDER"); v != "" {
		Cfg.Notification.Provider = v
	}

	// Notification - AWS SNS
	if v = os.Getenv("SC_NOTIFICATION_AWS_REGION"); v != "" {
		Cfg.Notification.AwsRegion = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_ACCESS_KEY_ID"); v != "" {
		Cfg.Notification.AwsAccessKeyId = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_SECRET_ACCESS_KEY"); v != "" {
		Cfg.Notification.AwsSecretAccessKey = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_APPLICATION_ARN_IOS"); v != "" {
		Cfg.Notification.AwsApplicationArnIos = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_APPLICATION_ARN_ANDROID"); v != "" {
		Cfg.Notification.AwsApplicationArnAndroid = v
	}

	// Realtime Server
	if v = os.Getenv("SC_REALTIME_SERVER_ENDPOINT"); v != "" {
		Cfg.RealtimeServer.Endpoint = v
	}
}

func parseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&Cfg.Port, "port", Cfg.Port, "")
	flag.BoolVar(&Cfg.Profiling, "profiling", Cfg.Profiling, "")
	flag.BoolVar(&Cfg.ErrorLogging, "errorLogging", Cfg.ErrorLogging, "")

	// Logging
	flag.StringVar(&Cfg.Logging.Level, "logging.level", Cfg.Logging.Level, "")

	// Storage
	flag.StringVar(&Cfg.Storage.Provider, "storage.provider", Cfg.Storage.Provider, "")
	flag.StringVar(&Cfg.Storage.UploadBucket, "storage.uploadBucket", Cfg.Storage.UploadBucket, "")
	flag.StringVar(&Cfg.Storage.UploadDirectory, "storage.uploadDirectory", Cfg.Storage.UploadDirectory, "")
	flag.StringVar(&Cfg.Storage.ThumbnailBucket, "storage.thumbnailBucket", Cfg.Storage.ThumbnailBucket, "")
	flag.StringVar(&Cfg.Storage.ThumbnailDirectory, "storage.thumbnailDirectory", Cfg.Storage.ThumbnailDirectory, "")

	// Storage - Local
	flag.StringVar(&Cfg.Storage.BaseUrl, "storage.baseUrl", Cfg.Storage.BaseUrl, "")
	flag.StringVar(&Cfg.Storage.LocalPath, "storage.localPath", Cfg.Storage.LocalPath, "")

	// Storage - GCP Storage
	flag.StringVar(&Cfg.Storage.GcpProjectId, "storage.gcpProjectId", Cfg.Storage.GcpProjectId, "")
	flag.StringVar(&Cfg.Storage.GcpJwtPath, "storage.gcpJwtPath", Cfg.Storage.GcpJwtPath, "")

	// Storage - AWS S3
	flag.StringVar(&Cfg.Storage.AwsRegion, "storage.awsRegion", Cfg.Storage.AwsRegion, "")
	flag.StringVar(&Cfg.Storage.AwsAccessKeyId, "storage.awsAccessKeyId", Cfg.Storage.AwsAccessKeyId, "")
	flag.StringVar(&Cfg.Storage.AwsSecretAccessKey, "storage.awsSecretAccessKey", Cfg.Storage.AwsSecretAccessKey, "")

	// Datastore
	flag.StringVar(&Cfg.Datastore.Provider, "datastore.provider", Cfg.Datastore.Provider, "")
	flag.StringVar(&Cfg.Datastore.TableNamePrefix, "datastore.tableNamePrefix", Cfg.Datastore.TableNamePrefix, "")

	// Datastore - SQLite
	flag.StringVar(&Cfg.Datastore.SqlitePath, "datastore.sqlitePath", Cfg.Datastore.SqlitePath, "")

	// Datastore - MySQL, GCP SQL
	flag.StringVar(&Cfg.Datastore.User, "datastore.user", Cfg.Datastore.User, "")
	flag.StringVar(&Cfg.Datastore.Password, "datastore.password", Cfg.Datastore.Password, "")
	flag.StringVar(&Cfg.Datastore.Database, "datastore.database", Cfg.Datastore.Database, "")
	flag.StringVar(&Cfg.Datastore.MasterHost, "datastore.masterHost", Cfg.Datastore.MasterHost, "")
	flag.StringVar(&Cfg.Datastore.MasterPort, "datastore.masterPort", Cfg.Datastore.MasterPort, "")
	flag.StringVar(&Cfg.Datastore.MaxIdleConnection, "datastore.maxIdleConnection", Cfg.Datastore.MaxIdleConnection, "")
	flag.StringVar(&Cfg.Datastore.MaxOpenConnection, "datastore.maxOpenConnection", Cfg.Datastore.MaxOpenConnection, "")
	flag.StringVar(&Cfg.Datastore.UseSSL, "datastore.useSSL", Cfg.Datastore.UseSSL, "")
	flag.StringVar(&Cfg.Datastore.ServerName, "datastore.serverName", Cfg.Datastore.ServerName, "")
	flag.StringVar(&Cfg.Datastore.ServerCaPath, "datastore.serverCaPath", Cfg.Datastore.ServerCaPath, "")
	flag.StringVar(&Cfg.Datastore.ClientCertPath, "datastore.clientCertPath", Cfg.Datastore.ClientCertPath, "")
	flag.StringVar(&Cfg.Datastore.ClientKeyPath, "datastore.clientKeyPath", Cfg.Datastore.ClientKeyPath, "")

	// Datastore -GCP SQL
	flag.StringVar(&Cfg.Datastore.GcpProjectId, "datastore.gcpProjectId", Cfg.Datastore.GcpProjectId, "")

	// Messaging
	flag.StringVar(&Cfg.Messaging.Provider, "messaging.provider", Cfg.Messaging.Provider, "")
	flag.StringVar(&Cfg.Messaging.ThumbnailTopic, "messaging.thumbnailTopic", Cfg.Messaging.ThumbnailTopic, "")

	// Messaging - GCP Pubsub
	flag.StringVar(&Cfg.Messaging.GcpProjectId, "messaging.gcpProjectId", Cfg.Messaging.GcpProjectId, "")
	flag.StringVar(&Cfg.Messaging.GcpJwtPath, "messaging.gcpJwtPath", Cfg.Messaging.GcpJwtPath, "")

	// Notification
	flag.StringVar(&Cfg.Notification.Provider, "notification.provider", Cfg.Notification.Provider, "")
	flag.StringVar(&Cfg.Notification.RoomTopicNamePrefix, "notification.roomTopicNamePrefix", Cfg.Notification.RoomTopicNamePrefix, "")

	// Notification - AWS SNS
	flag.StringVar(&Cfg.Notification.AwsRegion, "notification.awsRegion", Cfg.Notification.AwsRegion, "")
	flag.StringVar(&Cfg.Notification.AwsAccessKeyId, "notification.awsAccessKeyId", Cfg.Notification.AwsAccessKeyId, "")
	flag.StringVar(&Cfg.Notification.AwsSecretAccessKey, "notification.awsSecretAccessKey", Cfg.Notification.AwsSecretAccessKey, "")
	flag.StringVar(&Cfg.Notification.AwsApplicationArnIos, "notification.awsApplicationArnIos", Cfg.Notification.AwsApplicationArnIos, "")
	flag.StringVar(&Cfg.Notification.AwsApplicationArnAndroid, "notification.awsApplicationArnAndroid", Cfg.Notification.AwsApplicationArnAndroid, "")

	// Realtime Server
	flag.StringVar(&Cfg.RealtimeServer.Endpoint, "realtimeServer.endpoint", Cfg.RealtimeServer.Endpoint, "")

	flag.Parse()
}
