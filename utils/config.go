package utils

import (
	"flag"
	"io/ioutil"
	"os"

	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	APP_NAME      = "swagchat-api"
	API_VERSION   = "v0"
	BUILD_VERSION = "v0.9.1"

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
	Version      string
	Port         string
	Profiling    bool
	DemoPage     bool `yaml:"demoPage"`
	ErrorLogging bool `yaml:"errorLogging"`
	Logging      *Logging
	Storage      *Storage
	Datastore    *Datastore
	Rtm          *Rtm
	Notification *Notification
}

type Logging struct {
	Level string
}

type Storage struct {
	Provider string

	// Local
	BaseUrl   string `yaml:"baseUrl"`
	LocalPath string `yaml:"localPath"`

	// GCP Storage, AWS S3
	UploadBucket       string `yaml:"uploadBucket"`
	UploadDirectory    string `yaml:"uploadDirectory"`
	ThumbnailBucket    string `yaml:"thumbnailBucket"`
	ThumbnailDirectory string `yaml:"thumbnailDirectory"`

	// GCP Storage
	GcpProjectId string `yaml:"gcpProjectId"`
	GcpJwtPath   string `yaml:"gcpJwtPath"`

	// AWS S3
	AwsRegion          string `yaml:"awsRegion"`
	AwsAccessKeyId     string `yaml:"awsAccessKeyId"`
	AwsSecretAccessKey string `yaml:"awsSecretAccessKey"`
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
	MaxIdleConnection string `yaml:"maxIdleConnection"`
	MaxOpenConnection string `yaml:"maxOpenConnection"`
	Master            *ServerInfo
	Replicas          []*ServerInfo
}

type ServerInfo struct {
	Host           string
	Port           string
	ServerName     string `yaml:"serverName"`
	ServerCaPath   string `yaml:"serverCaPath"`
	ClientCertPath string `yaml:"clientCertPath"`
	ClientKeyPath  string `yaml:"clientKeyPath"`
}

type Rtm struct {
	Provider       string
	DirectEndpoint string `yaml:"directEndpoint"`
	QueEndpoint    string `yaml:"queEndpoint"`
	QueTopic       string `yaml:"queTopic"`
}

type Notification struct {
	Provider            string
	RoomTopicNamePrefix string `yaml:"roomTopicNamePrefix"`
	DefaultBadgeCount   string `yaml:"defaultBadgeCount"`

	// AWS SNS
	AwsRegion                string `yaml:"awsRegion"`
	AwsAccessKeyId           string `yaml:"awsAccessKeyId"`
	AwsSecretAccessKey       string `yaml:"awsSecretAccessKey"`
	AwsApplicationArnIos     string `yaml:"awsApplicationArnIos"`
	AwsApplicationArnAndroid string `yaml:"awsApplicationArnAndroid"`
}

func setupConfig() {
	loadDefaultSettings()
	loadYaml()
	loadEnvironment()
	parseFlag()
}

func loadDefaultSettings() {
	port := "9000"

	logging := &Logging{
		Level: "development",
	}

	storage := &Storage{
		Provider:  "local",
		BaseUrl:   AppendStrings("/", API_VERSION, "/assets"),
		LocalPath: "data/assets",
	}

	datastore := &Datastore{
		Provider:          "sqlite",
		SqlitePath:        "/tmp/swagchat.db",
		MaxIdleConnection: "10",
		MaxOpenConnection: "10",
	}

	rtm := &Rtm{
		Provider:       "",
		DirectEndpoint: "",
		QueEndpoint:    "",
		QueTopic:       "",
	}

	notification := &Notification{}

	Cfg = &Config{
		Version:      "0",
		Port:         port,
		Profiling:    false,
		DemoPage:     false,
		ErrorLogging: false,
		Logging:      logging,
		Storage:      storage,
		Datastore:    datastore,
		Rtm:          rtm,
		Notification: notification,
	}
}

func loadYaml() {
	buf, _ := ioutil.ReadFile("config/swagchat.yaml")
	yaml.Unmarshal(buf, Cfg)
}

func loadEnvironment() {
	var v string

	if v = os.Getenv("PORT"); v != "" {
		Cfg.Port = v
	}
	if v = os.Getenv("SC_PORT"); v != "" {
		Cfg.Port = v
	}
	if v = os.Getenv("SC_PROFILING"); v != "" {
		if v == "true" {
			Cfg.Profiling = true
		} else if v == "false" {
			Cfg.Profiling = false
		}
	}
	if v = os.Getenv("SC_DEMO_PAGE"); v != "" {
		if v == "true" {
			Cfg.DemoPage = true
		} else if v == "false" {
			Cfg.DemoPage = false
		}
	}
	if v = os.Getenv("SC_ERROR_LOGGING"); v != "" {
		if v == "true" {
			Cfg.ErrorLogging = true
		} else if v == "false" {
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

	// Storage - Local
	if v = os.Getenv("SC_STORAGE_BASE_URL"); v != "" {
		Cfg.Storage.BaseUrl = v
	}
	if v = os.Getenv("SC_STORAGE_LOCAL_PATH"); v != "" {
		Cfg.Storage.LocalPath = v
	}

	// Storage - GCP Storage, AWS S3
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
	if v = os.Getenv("SC_DATASTORE_MAX_IDLE_CONNECTION"); v != "" {
		Cfg.Datastore.MaxIdleConnection = v
	}
	if v = os.Getenv("SC_DATASTORE_MAX_OPEN_CONNECTION"); v != "" {
		Cfg.Datastore.MaxOpenConnection = v
	}

	var master *ServerInfo
	mHost := os.Getenv("SC_DATASTORE_MASTER_HOST")
	mPort := os.Getenv("SC_DATASTORE_MASTER_PORT")
	if mHost != "" && mPort != "" {
		master = &ServerInfo{}
		master.Host = mHost
		master.Port = mPort
		Cfg.Datastore.Master = master
		mServerName := os.Getenv("SC_DATASTORE_MASTER_SERVER_NAME")
		mServerCaPath := os.Getenv("SC_DATASTORE_MASTER_SERVER_CA_PATH")
		mClientCertPath := os.Getenv("SC_DATASTORE_MASTER_CLIENT_CERT_PATH")
		mClientKeyPath := os.Getenv("SC_DATASTORE_MASTER_CLIENT_KEY_PATH")
		if mServerName != "" && mServerCaPath != "" && mClientCertPath != "" && mClientKeyPath != "" {
			master.ServerName = mServerName
			master.ServerCaPath = mServerCaPath
			master.ClientCertPath = mClientCertPath
			master.ClientKeyPath = mClientKeyPath
		}
	}

	var (
		rHosts          []string
		rPorts          []string
		rServerName     []string
		rServerCaPath   []string
		rClientCertPath []string
		rClientKeyPath  []string
	)
	if v = os.Getenv("SC_DATASTORE_REPLICA_HOSTS"); v != "" {
		rHosts = strings.Split(v, ",")
	}
	if v = os.Getenv("SC_DATASTORE_REPLICA_PORTS"); v != "" {
		rPorts = strings.Split(v, ",")
	}
	if v = os.Getenv("SC_DATASTORE_REPLICA_SERVER_NAMES"); v != "" {
		rServerName = strings.Split(v, ",")
	}
	if v = os.Getenv("SC_DATASTORE_REPLICA_SERVER_CA_PATHS"); v != "" {
		rServerCaPath = strings.Split(v, ",")
	}
	if v = os.Getenv("SC_DATASTORE_REPLICA_CLIENT_CERT_PATHS"); v != "" {
		rClientCertPath = strings.Split(v, ",")
	}
	if v = os.Getenv("SC_DATASTORE_REPLICA_CLIENT_KEY_PATHS"); v != "" {
		rClientKeyPath = strings.Split(v, ",")
	}
	if rHosts != nil && len(rHosts) != 0 && rPorts != nil && len(rPorts) != 0 && len(rHosts) == len(rPorts) {
		replicas := []*ServerInfo{}
		for i := 0; i < len(rHosts); i++ {
			replica := &ServerInfo{
				Host: rHosts[i],
				Port: rPorts[i],
			}
			replicas = append(replicas, replica)
		}
		Cfg.Datastore.Replicas = replicas

		if rServerName != nil && len(rServerName) != 0 && rServerCaPath != nil && len(rServerCaPath) != 0 && rClientCertPath != nil && len(rClientCertPath) != 0 && rClientKeyPath != nil && len(rClientKeyPath) != 0 &&
			(len(rHosts) == len(rServerName) && len(rHosts) == len(rServerCaPath) && len(rHosts) == len(rClientCertPath) && len(rHosts) == len(rClientKeyPath)) {
			for i := 0; i < len(rHosts); i++ {
				Cfg.Datastore.Replicas[i].ServerName = rServerName[i]
				Cfg.Datastore.Replicas[i].ServerCaPath = rServerCaPath[i]
				Cfg.Datastore.Replicas[i].ClientCertPath = rClientCertPath[i]
				Cfg.Datastore.Replicas[i].ClientKeyPath = rClientKeyPath[i]
			}
		}
	}

	// Rtm
	if v = os.Getenv("SC_RTM_PROVIDER"); v != "" {
		Cfg.Rtm.Provider = v
	}
	if v = os.Getenv("SC_RTM_DIRECT_ENDPOINT"); v != "" {
		Cfg.Rtm.DirectEndpoint = v
	}
	if v = os.Getenv("SC_RTM_QUE_ENDPOINT"); v != "" {
		Cfg.Rtm.QueEndpoint = v
	}
	if v = os.Getenv("SC_RTM_QUE_TOPIC"); v != "" {
		Cfg.Rtm.QueTopic = v
	}

	// Notification
	if v = os.Getenv("SC_NOTIFICATION_PROVIDER"); v != "" {
		Cfg.Notification.Provider = v
	}
	if v = os.Getenv("SC_NOTIFICATION_ROOM_TOPIC_NAME_PREFIX"); v != "" {
		Cfg.Notification.RoomTopicNamePrefix = v
	}
	if v = os.Getenv("SC_NOTIFICATION_DEFAULT_BADGE_COUNT"); v != "" {
		Cfg.Notification.DefaultBadgeCount = v
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
	if v = os.Getenv("SC_NOTIFICATION_AWS_APPLICATION_ARN_IOS"); v != "" {
		Cfg.Notification.AwsApplicationArnIos = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_APPLICATION_ARN_ANDROID"); v != "" {
		Cfg.Notification.AwsApplicationArnAndroid = v
	}
}

func parseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&Cfg.Port, "port", Cfg.Port, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

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
	flag.StringVar(&Cfg.Datastore.MaxIdleConnection, "datastore.maxIdleConnection", Cfg.Datastore.MaxIdleConnection, "")
	flag.StringVar(&Cfg.Datastore.MaxOpenConnection, "datastore.maxOpenConnection", Cfg.Datastore.MaxOpenConnection, "")

	var (
		mHostStr           string
		mPortsStr          string
		mServerNameStr     string
		mServerCaPathStr   string
		mClientCertPathStr string
		mClientKeyPathStr  string
	)
	flag.StringVar(&mHostStr, "datastore.masterHost", mHostStr, "")
	flag.StringVar(&mPortsStr, "datastore.masterPort", mPortsStr, "")
	if mHostStr != "" && mPortsStr != "" {
		Cfg.Datastore.Master = &ServerInfo{
			Host: mHostStr,
			Port: mPortsStr,
		}
		flag.StringVar(&mServerNameStr, "datastore.masterServerName", mServerNameStr, "")
		flag.StringVar(&mServerCaPathStr, "datastore.masterServerCaPath", mServerCaPathStr, "")
		flag.StringVar(&mClientCertPathStr, "datastore.masterClientCertPath", mClientCertPathStr, "")
		flag.StringVar(&mClientKeyPathStr, "datastore.masterClientKeyPath", mClientKeyPathStr, "")
		if mServerNameStr != "" && mServerCaPathStr != "" && mClientCertPathStr != "" && mClientKeyPathStr != "" {
			Cfg.Datastore.Master.ServerName = mServerNameStr
			Cfg.Datastore.Master.ServerCaPath = mServerCaPathStr
			Cfg.Datastore.Master.ClientCertPath = mClientCertPathStr
			Cfg.Datastore.Master.ClientKeyPath = mClientKeyPathStr
		}
	}

	var (
		rHostsStr           string
		rPortsStr           string
		rServerNamesStr     string
		rServerCaPathsStr   string
		rClientCertPathsStr string
		rClientKeyPathsStr  string
		rHosts              []string
		rPorts              []string
		rServerNames        []string
		rServerCaPaths      []string
		rClientCertPaths    []string
		rClientKeyPaths     []string
	)
	flag.StringVar(&rHostsStr, "datastore.replicaHosts", rHostsStr, "")
	flag.StringVar(&rPortsStr, "datastore.replicaPorts", rPortsStr, "")
	flag.StringVar(&rServerNamesStr, "datastore.replicaServerNames", rServerNamesStr, "")
	flag.StringVar(&rServerCaPathsStr, "datastore.replicaServerCaPaths", rServerCaPathsStr, "")
	flag.StringVar(&rClientCertPathsStr, "datastore.replicaClientCertPaths", rClientCertPathsStr, "")
	flag.StringVar(&rClientKeyPathsStr, "datastore.replicaClientKeyPaths", rClientKeyPathsStr, "")

	if rHostsStr != "" {
		rHosts = strings.Split(rHostsStr, ",")
	}
	if rPortsStr != "" {
		rPorts = strings.Split(rPortsStr, ",")
	}
	if rServerNamesStr != "" {
		rServerNames = strings.Split(rServerNamesStr, ",")
	}
	if rServerCaPathsStr != "" {
		rServerCaPaths = strings.Split(rServerCaPathsStr, ",")
	}
	if rClientCertPathsStr != "" {
		rClientCertPaths = strings.Split(rClientCertPathsStr, ",")
	}
	if rClientKeyPathsStr != "" {
		rClientKeyPaths = strings.Split(rClientKeyPathsStr, ",")
	}
	if rHosts != nil && len(rHosts) != 0 && rPorts != nil && len(rPorts) != 0 && len(rHosts) == len(rPorts) {
		replicas := []*ServerInfo{}
		for i := 0; i < len(rHosts); i++ {
			replica := &ServerInfo{
				Host: rHosts[i],
				Port: rPorts[i],
			}
			replicas = append(replicas, replica)
		}
		Cfg.Datastore.Replicas = replicas

		if rServerNames != nil && len(rServerNames) != 0 && rServerCaPaths != nil && len(rServerCaPaths) != 0 && rClientCertPaths != nil && len(rClientCertPaths) != 0 && rClientKeyPaths != nil && len(rClientKeyPaths) != 0 &&
			(len(rHosts) == len(rServerNames) && len(rHosts) == len(rServerCaPaths) && len(rHosts) == len(rClientCertPaths) && len(rHosts) == len(rClientKeyPaths)) {
			for i := 0; i < len(rHosts); i++ {
				Cfg.Datastore.Replicas[i].ServerName = rServerNames[i]
				Cfg.Datastore.Replicas[i].ServerCaPath = rServerCaPaths[i]
				Cfg.Datastore.Replicas[i].ClientCertPath = rClientCertPaths[i]
				Cfg.Datastore.Replicas[i].ClientKeyPath = rClientKeyPaths[i]
			}
		}
	}

	// Rtm
	flag.StringVar(&Cfg.Rtm.Provider, "realtimeMessaging.provider", Cfg.Rtm.Provider, "")
	flag.StringVar(&Cfg.Rtm.DirectEndpoint, "realtimeMessaging.directEndpoint", Cfg.Rtm.DirectEndpoint, "")
	flag.StringVar(&Cfg.Rtm.QueEndpoint, "realtimeMessaging.queEndpoint", Cfg.Rtm.QueEndpoint, "")
	flag.StringVar(&Cfg.Rtm.QueTopic, "realtimeMessaging.queTopic", Cfg.Rtm.QueTopic, "")

	// Notification
	flag.StringVar(&Cfg.Notification.Provider, "notification.provider", Cfg.Notification.Provider, "")
	flag.StringVar(&Cfg.Notification.RoomTopicNamePrefix, "notification.roomTopicNamePrefix", Cfg.Notification.RoomTopicNamePrefix, "")

	// Notification - AWS SNS
	flag.StringVar(&Cfg.Notification.AwsRegion, "notification.awsRegion", Cfg.Notification.AwsRegion, "")
	flag.StringVar(&Cfg.Notification.AwsAccessKeyId, "notification.awsAccessKeyId", Cfg.Notification.AwsAccessKeyId, "")
	flag.StringVar(&Cfg.Notification.AwsSecretAccessKey, "notification.awsSecretAccessKey", Cfg.Notification.AwsSecretAccessKey, "")
	flag.StringVar(&Cfg.Notification.AwsApplicationArnIos, "notification.awsApplicationArnIos", Cfg.Notification.AwsApplicationArnIos, "")
	flag.StringVar(&Cfg.Notification.AwsApplicationArnAndroid, "notification.awsApplicationArnAndroid", Cfg.Notification.AwsApplicationArnAndroid, "")
	flag.Parse()

	if profiling == "true" {
		Cfg.Profiling = true
	} else if profiling == "false" {
		Cfg.Profiling = false
	}

	if demoPage == "true" {
		Cfg.DemoPage = true
	} else if demoPage == "false" {
		Cfg.DemoPage = false
	}

	if errorLogging == "true" {
		Cfg.ErrorLogging = true
	} else if errorLogging == "false" {
		Cfg.ErrorLogging = false
	}
}
