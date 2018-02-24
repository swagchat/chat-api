package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// AppName is Application name
	AppName = "chat-api"
	// APIVersion is API version
	APIVersion = "v0"
	// BuildVersion is API build version
	BuildVersion = "v0.9.1"

	KeyLength       = 32
	TokenLength     = 32
	HeaderAPIKey    = "X-SwagChat-Api-Key"
	HeaderAPISecret = "X-SwagChat-Api-Secret"
	HeaderUserId    = "X-SwagChat-User-Id"
)

var (
	cfg           *config = NewConfig()
	IsShowVersion bool
)

type config struct {
	Version      string
	HttpPort     string `yaml:"httpPort"`
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

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	logging := &Logging{
		Level: "development",
	}

	storage := &Storage{
		Provider:  "local",
		BaseUrl:   "/assets",
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

	c := &config{
		Version:      "0",
		HttpPort:     "8101",
		Profiling:    false,
		DemoPage:     false,
		ErrorLogging: false,
		Logging:      logging,
		Storage:      storage,
		Datastore:    datastore,
		Rtm:          rtm,
		Notification: notification,
	}

	c.LoadYaml()
	c.LoadEnvironment()
	c.ParseFlag()
	c.PrintConfig()

	return c
}

func GetConfig() *config {
	return cfg
}

func (c *config) LoadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) LoadEnvironment() {
	var v string

	if v = os.Getenv("HTTP_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("SC_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("SC_PROFILING"); v != "" {
		if v == "true" {
			c.Profiling = true
		} else if v == "false" {
			c.Profiling = false
		}
	}
	if v = os.Getenv("SC_DEMO_PAGE"); v != "" {
		if v == "true" {
			c.DemoPage = true
		} else if v == "false" {
			c.DemoPage = false
		}
	}
	if v = os.Getenv("SC_ERROR_LOGGING"); v != "" {
		if v == "true" {
			c.ErrorLogging = true
		} else if v == "false" {
			c.ErrorLogging = false
		}
	}

	// Logging
	if v = os.Getenv("SC_LOGGING_LEVEL"); v != "" {
		c.Logging.Level = v
	}

	// Storage
	if v = os.Getenv("SC_STORAGE_PROVIDER"); v != "" {
		c.Storage.Provider = v
	}

	// Storage - Local
	if v = os.Getenv("SC_STORAGE_BASE_URL"); v != "" {
		c.Storage.BaseUrl = v
	}
	if v = os.Getenv("SC_STORAGE_LOCAL_PATH"); v != "" {
		c.Storage.LocalPath = v
	}

	// Storage - GCP Storage, AWS S3
	if v = os.Getenv("SC_STORAGE_UPLOAD_BUCKET"); v != "" {
		c.Storage.UploadBucket = v
	}
	if v = os.Getenv("SC_STORAGE_UPLOAD_DIRECTORY"); v != "" {
		c.Storage.UploadDirectory = v
	}
	if v = os.Getenv("SC_STORAGE_THUMBNAIL_BUCKET"); v != "" {
		c.Storage.ThumbnailBucket = v
	}
	if v = os.Getenv("SC_STORAGE_THUMBNAIL_DIRECTORY"); v != "" {
		c.Storage.ThumbnailDirectory = v
	}

	// Storage - GCP Storage
	if v = os.Getenv("SC_STORAGE_GCP_PROJECT_ID"); v != "" {
		c.Storage.GcpProjectId = v
	}
	if v = os.Getenv("SC_STORAGE_GCP_JWT_PATH"); v != "" {
		c.Storage.GcpJwtPath = v
	}

	// Storage - AWS S3
	if v = os.Getenv("SC_STORAGE_AWS_REGION"); v != "" {
		c.Storage.AwsRegion = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_ACCESS_KEY_ID"); v != "" {
		c.Storage.AwsAccessKeyId = v
	}
	if v = os.Getenv("SC_STORAGE_AWS_SECRET_ACCESS_KEY"); v != "" {
		c.Storage.AwsSecretAccessKey = v
	}

	// Datastore
	if v = os.Getenv("SC_DATASTORE_PROVIDER"); v != "" {
		c.Datastore.Provider = v
	}
	if v = os.Getenv("SC_DATASTORE_TABLE_NAME_PREFIX"); v != "" {
		c.Datastore.TableNamePrefix = v
	}

	// Datastore - SQLite
	if v = os.Getenv("SC_DATASTORE_SQLITE_PATH"); v != "" {
		c.Datastore.SqlitePath = v
	}

	// Datastore - MySQL, GCP SQL
	if v = os.Getenv("SC_DATASTORE_USER"); v != "" {
		c.Datastore.User = v
	}
	if v = os.Getenv("SC_DATASTORE_PASSWORD"); v != "" {
		c.Datastore.Password = v
	}
	if v = os.Getenv("SC_DATASTORE_DATABASE"); v != "" {
		c.Datastore.Database = v
	}
	if v = os.Getenv("SC_DATASTORE_MAX_IDLE_CONNECTION"); v != "" {
		c.Datastore.MaxIdleConnection = v
	}
	if v = os.Getenv("SC_DATASTORE_MAX_OPEN_CONNECTION"); v != "" {
		c.Datastore.MaxOpenConnection = v
	}

	var master *ServerInfo
	mHost := os.Getenv("SC_DATASTORE_MASTER_HOST")
	mPort := os.Getenv("SC_DATASTORE_MASTER_PORT")
	if mHost != "" && mPort != "" {
		master = &ServerInfo{}
		master.Host = mHost
		master.Port = mPort
		c.Datastore.Master = master
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
		c.Datastore.Replicas = replicas

		if rServerName != nil && len(rServerName) != 0 && rServerCaPath != nil && len(rServerCaPath) != 0 && rClientCertPath != nil && len(rClientCertPath) != 0 && rClientKeyPath != nil && len(rClientKeyPath) != 0 &&
			(len(rHosts) == len(rServerName) && len(rHosts) == len(rServerCaPath) && len(rHosts) == len(rClientCertPath) && len(rHosts) == len(rClientKeyPath)) {
			for i := 0; i < len(rHosts); i++ {
				c.Datastore.Replicas[i].ServerName = rServerName[i]
				c.Datastore.Replicas[i].ServerCaPath = rServerCaPath[i]
				c.Datastore.Replicas[i].ClientCertPath = rClientCertPath[i]
				c.Datastore.Replicas[i].ClientKeyPath = rClientKeyPath[i]
			}
		}
	}

	// Rtm
	if v = os.Getenv("SC_RTM_PROVIDER"); v != "" {
		c.Rtm.Provider = v
	}
	if v = os.Getenv("SC_RTM_DIRECT_ENDPOINT"); v != "" {
		c.Rtm.DirectEndpoint = v
	}
	if v = os.Getenv("SC_RTM_QUE_ENDPOINT"); v != "" {
		c.Rtm.QueEndpoint = v
	}
	if v = os.Getenv("SC_RTM_QUE_TOPIC"); v != "" {
		c.Rtm.QueTopic = v
	}

	// Notification
	if v = os.Getenv("SC_NOTIFICATION_PROVIDER"); v != "" {
		c.Notification.Provider = v
	}
	if v = os.Getenv("SC_NOTIFICATION_ROOM_TOPIC_NAME_PREFIX"); v != "" {
		c.Notification.RoomTopicNamePrefix = v
	}
	if v = os.Getenv("SC_NOTIFICATION_DEFAULT_BADGE_COUNT"); v != "" {
		c.Notification.DefaultBadgeCount = v
	}

	// Notification - AWS SNS
	if v = os.Getenv("SC_NOTIFICATION_AWS_REGION"); v != "" {
		c.Notification.AwsRegion = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_ACCESS_KEY_ID"); v != "" {
		c.Notification.AwsAccessKeyId = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_SECRET_ACCESS_KEY"); v != "" {
		c.Notification.AwsSecretAccessKey = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_APPLICATION_ARN_IOS"); v != "" {
		c.Notification.AwsApplicationArnIos = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AWS_APPLICATION_ARN_ANDROID"); v != "" {
		c.Notification.AwsApplicationArnAndroid = v
	}
}

func (c *config) ParseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HttpPort, "httpPort", c.HttpPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	// Storage
	flag.StringVar(&c.Storage.Provider, "storage.provider", c.Storage.Provider, "")
	flag.StringVar(&c.Storage.UploadBucket, "storage.uploadBucket", c.Storage.UploadBucket, "")
	flag.StringVar(&c.Storage.UploadDirectory, "storage.uploadDirectory", c.Storage.UploadDirectory, "")
	flag.StringVar(&c.Storage.ThumbnailBucket, "storage.thumbnailBucket", c.Storage.ThumbnailBucket, "")
	flag.StringVar(&c.Storage.ThumbnailDirectory, "storage.thumbnailDirectory", c.Storage.ThumbnailDirectory, "")

	// Storage - Local
	flag.StringVar(&c.Storage.BaseUrl, "storage.baseUrl", c.Storage.BaseUrl, "")
	flag.StringVar(&c.Storage.LocalPath, "storage.localPath", c.Storage.LocalPath, "")

	// Storage - GCP Storage
	flag.StringVar(&c.Storage.GcpProjectId, "storage.gcpProjectId", c.Storage.GcpProjectId, "")
	flag.StringVar(&c.Storage.GcpJwtPath, "storage.gcpJwtPath", c.Storage.GcpJwtPath, "")

	// Storage - AWS S3
	flag.StringVar(&c.Storage.AwsRegion, "storage.awsRegion", c.Storage.AwsRegion, "")
	flag.StringVar(&c.Storage.AwsAccessKeyId, "storage.awsAccessKeyId", c.Storage.AwsAccessKeyId, "")
	flag.StringVar(&c.Storage.AwsSecretAccessKey, "storage.awsSecretAccessKey", c.Storage.AwsSecretAccessKey, "")

	// Datastore
	flag.StringVar(&c.Datastore.Provider, "datastore.provider", c.Datastore.Provider, "")
	flag.StringVar(&c.Datastore.TableNamePrefix, "datastore.tableNamePrefix", c.Datastore.TableNamePrefix, "")

	// Datastore - SQLite
	flag.StringVar(&c.Datastore.SqlitePath, "datastore.sqlitePath", c.Datastore.SqlitePath, "")

	// Datastore - MySQL, GCP SQL
	flag.StringVar(&c.Datastore.User, "datastore.user", c.Datastore.User, "")
	flag.StringVar(&c.Datastore.Password, "datastore.password", c.Datastore.Password, "")
	flag.StringVar(&c.Datastore.Database, "datastore.database", c.Datastore.Database, "")
	flag.StringVar(&c.Datastore.MaxIdleConnection, "datastore.maxIdleConnection", c.Datastore.MaxIdleConnection, "")
	flag.StringVar(&c.Datastore.MaxOpenConnection, "datastore.maxOpenConnection", c.Datastore.MaxOpenConnection, "")

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
		c.Datastore.Master = &ServerInfo{
			Host: mHostStr,
			Port: mPortsStr,
		}
		flag.StringVar(&mServerNameStr, "datastore.masterServerName", mServerNameStr, "")
		flag.StringVar(&mServerCaPathStr, "datastore.masterServerCaPath", mServerCaPathStr, "")
		flag.StringVar(&mClientCertPathStr, "datastore.masterClientCertPath", mClientCertPathStr, "")
		flag.StringVar(&mClientKeyPathStr, "datastore.masterClientKeyPath", mClientKeyPathStr, "")
		if mServerNameStr != "" && mServerCaPathStr != "" && mClientCertPathStr != "" && mClientKeyPathStr != "" {
			c.Datastore.Master.ServerName = mServerNameStr
			c.Datastore.Master.ServerCaPath = mServerCaPathStr
			c.Datastore.Master.ClientCertPath = mClientCertPathStr
			c.Datastore.Master.ClientKeyPath = mClientKeyPathStr
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
		c.Datastore.Replicas = replicas

		if rServerNames != nil && len(rServerNames) != 0 && rServerCaPaths != nil && len(rServerCaPaths) != 0 && rClientCertPaths != nil && len(rClientCertPaths) != 0 && rClientKeyPaths != nil && len(rClientKeyPaths) != 0 &&
			(len(rHosts) == len(rServerNames) && len(rHosts) == len(rServerCaPaths) && len(rHosts) == len(rClientCertPaths) && len(rHosts) == len(rClientKeyPaths)) {
			for i := 0; i < len(rHosts); i++ {
				c.Datastore.Replicas[i].ServerName = rServerNames[i]
				c.Datastore.Replicas[i].ServerCaPath = rServerCaPaths[i]
				c.Datastore.Replicas[i].ClientCertPath = rClientCertPaths[i]
				c.Datastore.Replicas[i].ClientKeyPath = rClientKeyPaths[i]
			}
		}
	}

	// Rtm
	flag.StringVar(&c.Rtm.Provider, "realtimeMessaging.provider", c.Rtm.Provider, "")
	flag.StringVar(&c.Rtm.DirectEndpoint, "realtimeMessaging.directEndpoint", c.Rtm.DirectEndpoint, "")
	flag.StringVar(&c.Rtm.QueEndpoint, "realtimeMessaging.queEndpoint", c.Rtm.QueEndpoint, "")
	flag.StringVar(&c.Rtm.QueTopic, "realtimeMessaging.queTopic", c.Rtm.QueTopic, "")

	// Notification
	flag.StringVar(&c.Notification.Provider, "notification.provider", c.Notification.Provider, "")
	flag.StringVar(&c.Notification.RoomTopicNamePrefix, "notification.roomTopicNamePrefix", c.Notification.RoomTopicNamePrefix, "")

	// Notification - AWS SNS
	flag.StringVar(&c.Notification.AwsRegion, "notification.awsRegion", c.Notification.AwsRegion, "")
	flag.StringVar(&c.Notification.AwsAccessKeyId, "notification.awsAccessKeyId", c.Notification.AwsAccessKeyId, "")
	flag.StringVar(&c.Notification.AwsSecretAccessKey, "notification.awsSecretAccessKey", c.Notification.AwsSecretAccessKey, "")
	flag.StringVar(&c.Notification.AwsApplicationArnIos, "notification.awsApplicationArnIos", c.Notification.AwsApplicationArnIos, "")
	flag.StringVar(&c.Notification.AwsApplicationArnAndroid, "notification.awsApplicationArnAndroid", c.Notification.AwsApplicationArnAndroid, "")
	flag.Parse()

	if profiling == "true" {
		c.Profiling = true
	} else if profiling == "false" {
		c.Profiling = false
	}

	if demoPage == "true" {
		c.DemoPage = true
	} else if demoPage == "false" {
		c.DemoPage = false
	}

	if errorLogging == "true" {
		c.ErrorLogging = true
	} else if errorLogging == "false" {
		c.ErrorLogging = false
	}
}

func (c *config) PrintConfig() {
	fmt.Printf("%+v\n", c)
}
