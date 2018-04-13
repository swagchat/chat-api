package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"strings"

	yaml "gopkg.in/yaml.v2"
)

type key int

const (
	// AppName is Application name
	AppName = "chat-api"
	// APIVersion is API version
	APIVersion = "0"
	// BuildVersion is API build version
	BuildVersion = "0.9.1"

	// KeyLength is key length
	KeyLength = 32
	// TokenLength is token length
	TokenLength = 32

	// HeaderUserID is http header for userID
	HeaderUserID = "X-Sub"
	// HeaderUserName is http header for username
	HeaderUsername = "X-Preferred-Username"
	// HeaderRealm is http header for realm
	HeaderRealm = "X-Realm"

	CtxDsCfg key = iota
	CtxIsAppClient
	CtxUserID
	CtxRoomUser
	CtxSubscription
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
	Auth         *Auth
	Logging      *Logging
	Storage      *Storage
	Datastore    *Datastore
	RTM          *RTM
	Notification *Notification
}

type Auth struct {
	DefaultUsernameJWTClaimName string `yaml:"defaultUsernameJWTClaimName"`
}

type Logging struct {
	Level string
}

type Storage struct {
	Provider string

	Local struct {
		Path string `yaml:"localPath"`
	}

	GCS struct {
		ProjectID          string `yaml:"projectId"`
		JwtPath            string `yaml:"jwtPath"`
		UploadBucket       string `yaml:"uploadBucket"`
		UploadDirectory    string `yaml:"uploadDirectory"`
		ThumbnailBucket    string `yaml:"thumbnailBucket"`
		ThumbnailDirectory string `yaml:"thumbnailDirectory"`
	}

	AWSS3 struct {
		Region             string `yaml:"region"`
		AccessKeyID        string `yaml:"accessKeyId"`
		SecretAccessKey    string `yaml:"secretAccessKey"`
		UploadBucket       string `yaml:"uploadBucket"`
		UploadDirectory    string `yaml:"uploadDirectory"`
		ThumbnailBucket    string `yaml:"thumbnailBucket"`
		ThumbnailDirectory string `yaml:"thumbnailDirectory"`
	}
}

type Datastore struct {
	Dynamic  bool
	Provider string

	User              string
	Password          string
	Database          string
	TableNamePrefix   string `yaml:"tableNamePrefix"`
	MaxIdleConnection string `yaml:"maxIdleConnection"`
	MaxOpenConnection string `yaml:"maxOpenConnection"`
	Master            *ServerInfo
	Replicas          []*ServerInfo

	// SQLite
	SQLite struct {
		Path string `yaml:"path"`
	} `yaml:"sqlite"`
}

type ServerInfo struct {
	Host           string
	Port           string
	ServerName     string `yaml:"serverName"`
	ServerCaPath   string `yaml:"serverCaPath"`
	ClientCertPath string `yaml:"clientCertPath"`
	ClientKeyPath  string `yaml:"clientKeyPath"`
}

type RTM struct {
	Provider string

	Direct struct {
		Endpoint string `yaml:"directEndpoint"`
	}

	Kafka struct {
		Host    string
		Port    string
		GroupID string `yaml:"groupId"`
		Topic   string
	}

	NSQ struct {
		QueEndpoint string `yaml:"queEndpoint"`
		QueTopic    string `yaml:"queTopic"`
	}
}

type Notification struct {
	Provider            string
	RoomTopicNamePrefix string `yaml:"roomTopicNamePrefix"`
	DefaultBadgeCount   string `yaml:"defaultBadgeCount"`

	// Amazon SNS
	AmazonSNS struct {
		Region                string `yaml:"awsRegion"`
		AccessKeyID           string `yaml:"accessKeyId"`
		SecretAccessKey       string `yaml:"secretAccessKey"`
		ApplicationArnIos     string `yaml:"applicationArnIos"`
		ApplicationArnAndroid string `yaml:"applicationArnAndroid"`
	}
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	auth := &Auth{
		DefaultUsernameJWTClaimName: "preferred_username",
	}

	logging := &Logging{
		Level: "development",
	}

	storage := &Storage{
		Provider: "local",
	}
	storage.Local.Path = "data/assets"

	datastore := &Datastore{
		Dynamic:  false,
		Provider: "sqlite",
	}

	rtm := &RTM{}

	notification := &Notification{}

	c := &config{
		Version:      "0",
		HttpPort:     "8101",
		Profiling:    false,
		DemoPage:     false,
		ErrorLogging: false,
		Auth:         auth,
		Logging:      logging,
		Storage:      storage,
		Datastore:    datastore,
		RTM:          rtm,
		Notification: notification,
	}

	c.loadYaml()
	c.loadEnvironment()
	c.parseFlag()
	c.after()

	return c
}

// Config is get config
func Config() *config {
	return cfg
}

func (c *config) loadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) loadEnvironment() {
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

	// Auth
	if v = os.Getenv("SC_AUTH_DEFAULT_USERNAME_JWT_CLAIM_NAME"); v != "" {
		c.Auth.DefaultUsernameJWTClaimName = v
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
	if v = os.Getenv("SC_STORAGE_LOCAL_PATH"); v != "" {
		c.Storage.Local.Path = v
	}

	// Storage - Google Cloud Storage
	if v = os.Getenv("SC_STORAGE_GCS_PROJECT_ID"); v != "" {
		c.Storage.GCS.ProjectID = v
	}
	if v = os.Getenv("SC_STORAGE_GCS_JWT_PATH"); v != "" {
		c.Storage.GCS.JwtPath = v
	}
	if v = os.Getenv("SC_STORAGE_GCS_UPLOAD_BUCKET"); v != "" {
		c.Storage.GCS.UploadBucket = v
	}
	if v = os.Getenv("SC_STORAGE_GCS_UPLOAD_DIRECTORY"); v != "" {
		c.Storage.GCS.UploadDirectory = v
	}
	if v = os.Getenv("SC_STORAGE_GCS_THUMBNAIL_BUCKET"); v != "" {
		c.Storage.GCS.ThumbnailBucket = v
	}
	if v = os.Getenv("SC_STORAGE_GCS_THUMBNAIL_DIRECTORY"); v != "" {
		c.Storage.GCS.ThumbnailDirectory = v
	}

	// Storage - AWS S3
	if v = os.Getenv("SC_STORAGE_AWSS3_REGION"); v != "" {
		c.Storage.AWSS3.Region = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_ACCESS_KEY_ID"); v != "" {
		c.Storage.AWSS3.AccessKeyID = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_SECRET_ACCESS_KEY"); v != "" {
		c.Storage.AWSS3.SecretAccessKey = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_UPLOAD_BUCKET"); v != "" {
		c.Storage.AWSS3.UploadBucket = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_UPLOAD_DIRECTORY"); v != "" {
		c.Storage.AWSS3.UploadDirectory = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_THUMBNAIL_BUCKET"); v != "" {
		c.Storage.AWSS3.ThumbnailBucket = v
	}
	if v = os.Getenv("SC_STORAGE_AWSS3_THUMBNAIL_DIRECTORY"); v != "" {
		c.Storage.AWSS3.ThumbnailDirectory = v
	}

	// Datastore
	if v = os.Getenv("SC_DATASTORE_DYNAMIC"); v != "" {
		if v == "true" {
			c.Datastore.Dynamic = true
		} else if v == "false" {
			c.Datastore.Dynamic = false
		}
	}
	if v = os.Getenv("SC_DATASTORE_PROVIDER"); v != "" {
		c.Datastore.Provider = v
	}

	if v = os.Getenv("SC_DATASTORE_USER"); v != "" {
		c.Datastore.User = v
	}
	if v = os.Getenv("SC_DATASTORE_PASSWORD"); v != "" {
		c.Datastore.Password = v
	}
	if v = os.Getenv("SC_DATASTORE_DATABASE"); v != "" {
		c.Datastore.Database = v
	}
	if v = os.Getenv("SC_DATASTORE_TABLE_NAME_PREFIX"); v != "" {
		c.Datastore.TableNamePrefix = v
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

	// Datastore - SQLite
	if v = os.Getenv("SC_DATASTORE_SQLITE_PATH"); v != "" {
		c.Datastore.SQLite.Path = v
	}

	// RTM
	if v = os.Getenv("SC_RTM_PROVIDER"); v != "" {
		c.RTM.Provider = v
	}

	// RTM - Direct
	if v = os.Getenv("SC_RTM_DIRECT_ENDPOINT"); v != "" {
		c.RTM.Direct.Endpoint = v
	}

	// RTM - Kafka
	if v = os.Getenv("SC_RTM_KAFKA_HOST"); v != "" {
		c.RTM.Kafka.Host = v
	}
	if v = os.Getenv("SC_RTM_KAFKA_PORT"); v != "" {
		c.RTM.Kafka.Port = v
	}
	if v = os.Getenv("SC_RTM_KAFKA_GROUPID"); v != "" {
		c.RTM.Kafka.GroupID = v
	}
	if v = os.Getenv("SC_RTM_KAFKA_TOPIC"); v != "" {
		c.RTM.Kafka.Topic = v
	}

	// RTM - NSQ
	if v = os.Getenv("SC_RTM_NSQ_QUEENDPOINT"); v != "" {
		c.RTM.NSQ.QueEndpoint = v
	}
	if v = os.Getenv("SC_RTM_NSQ_QUETOPIC"); v != "" {
		c.RTM.NSQ.QueTopic = v
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

	// Notification - Amazon SNS
	if v = os.Getenv("SC_NOTIFICATION_AMAZONSNS_REGION"); v != "" {
		c.Notification.AmazonSNS.Region = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AMAZONSNS_ACCESS_KEY_ID"); v != "" {
		c.Notification.AmazonSNS.AccessKeyID = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AMAZONSNS_SECRET_ACCESS_KEY"); v != "" {
		c.Notification.AmazonSNS.SecretAccessKey = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AMAZONSNS_APPLICATION_ARN_IOS"); v != "" {
		c.Notification.AmazonSNS.ApplicationArnIos = v
	}
	if v = os.Getenv("SC_NOTIFICATION_AMAZONSNS_APPLICATION_ARN_ANDROID"); v != "" {
		c.Notification.AmazonSNS.ApplicationArnAndroid = v
	}
}

func (c *config) parseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HttpPort, "httpPort", c.HttpPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	// Auth
	flag.StringVar(&c.Auth.DefaultUsernameJWTClaimName, "auth.defaultUsernameJWTClaimName", c.Auth.DefaultUsernameJWTClaimName, "")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	// Storage
	flag.StringVar(&c.Storage.Provider, "storage.provider", c.Storage.Provider, "")

	// Storage - Local
	flag.StringVar(&c.Storage.Local.Path, "storage.local.path", c.Storage.Local.Path, "")

	// Storage - Google Cloud Storage
	flag.StringVar(&c.Storage.GCS.ProjectID, "storage.gcs.projectId", c.Storage.GCS.ProjectID, "")
	flag.StringVar(&c.Storage.GCS.JwtPath, "storage.gcs.jwtPath", c.Storage.GCS.JwtPath, "")
	flag.StringVar(&c.Storage.GCS.UploadBucket, "storage.gcs.uploadBucket", c.Storage.GCS.UploadBucket, "")
	flag.StringVar(&c.Storage.GCS.UploadDirectory, "storage.gcs.uploadDirectory", c.Storage.GCS.UploadDirectory, "")
	flag.StringVar(&c.Storage.GCS.ThumbnailBucket, "storage.gcs.thumbnailBucket", c.Storage.GCS.ThumbnailBucket, "")
	flag.StringVar(&c.Storage.GCS.ThumbnailDirectory, "storage.gcs.thumbnailDirectory", c.Storage.GCS.ThumbnailDirectory, "")

	// Storage - AWS S3
	flag.StringVar(&c.Storage.AWSS3.Region, "storage.awss3.region", c.Storage.AWSS3.Region, "")
	flag.StringVar(&c.Storage.AWSS3.AccessKeyID, "storage.awss3.accessKeyId", c.Storage.AWSS3.AccessKeyID, "")
	flag.StringVar(&c.Storage.AWSS3.SecretAccessKey, "storage.awss3.secretAccessKey", c.Storage.AWSS3.SecretAccessKey, "")
	flag.StringVar(&c.Storage.AWSS3.UploadBucket, "storage.awss3.uploadBucket", c.Storage.AWSS3.UploadBucket, "")
	flag.StringVar(&c.Storage.AWSS3.UploadDirectory, "storage.awss3.uploadDirectory", c.Storage.AWSS3.UploadDirectory, "")
	flag.StringVar(&c.Storage.AWSS3.ThumbnailBucket, "storage.awss3.thumbnailBucket", c.Storage.AWSS3.ThumbnailBucket, "")
	flag.StringVar(&c.Storage.AWSS3.ThumbnailDirectory, "storage.awss3.thumbnailDirectory", c.Storage.AWSS3.ThumbnailDirectory, "")

	// Datastore
	flag.BoolVar(&c.Datastore.Dynamic, "datastore.dynamic", c.Datastore.Dynamic, "")
	flag.StringVar(&c.Datastore.Provider, "datastore.provider", c.Datastore.Provider, "")
	flag.StringVar(&c.Datastore.TableNamePrefix, "datastore.tableNamePrefix", c.Datastore.TableNamePrefix, "")
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

	// Datastore - SQLite
	flag.StringVar(&c.Datastore.SQLite.Path, "datastore.sqlite.path", c.Datastore.SQLite.Path, "")

	// RTM
	flag.StringVar(&c.RTM.Provider, "rtm.provider", c.RTM.Provider, "")

	// RTM - Direct
	flag.StringVar(&c.RTM.Direct.Endpoint, "rtm.direct.endpoint", c.RTM.Direct.Endpoint, "")

	// RTM - kafka
	flag.StringVar(&c.RTM.Kafka.Host, "rtm.kafka.host", c.RTM.Kafka.Host, "")
	flag.StringVar(&c.RTM.Kafka.Port, "rtm.kafka.port", c.RTM.Kafka.Port, "")
	flag.StringVar(&c.RTM.Kafka.GroupID, "rtm.kafka.groupId", c.RTM.Kafka.GroupID, "")
	flag.StringVar(&c.RTM.Kafka.Topic, "rtm.kafka.topic", c.RTM.Kafka.Topic, "")

	// RTM - NSQ
	flag.StringVar(&c.RTM.NSQ.QueEndpoint, "rtm.nsq.queEndpoint", c.RTM.NSQ.QueEndpoint, "")
	flag.StringVar(&c.RTM.NSQ.QueTopic, "rtm.nsq.queTopic", c.RTM.NSQ.QueTopic, "")

	// Notification
	flag.StringVar(&c.Notification.Provider, "notification.provider", c.Notification.Provider, "")
	flag.StringVar(&c.Notification.RoomTopicNamePrefix, "notification.roomTopicNamePrefix", c.Notification.RoomTopicNamePrefix, "")

	// Notification - Amazon SNS
	flag.StringVar(&c.Notification.AmazonSNS.Region, "notification.amazonsns.region", c.Notification.AmazonSNS.Region, "")
	flag.StringVar(&c.Notification.AmazonSNS.AccessKeyID, "notification.amazonsns.accessKeyId", c.Notification.AmazonSNS.AccessKeyID, "")
	flag.StringVar(&c.Notification.AmazonSNS.SecretAccessKey, "notification.amazonsns.secretAccessKey", c.Notification.AmazonSNS.SecretAccessKey, "")
	flag.StringVar(&c.Notification.AmazonSNS.ApplicationArnIos, "notification.amazonsns.applicationArnIos", c.Notification.AmazonSNS.ApplicationArnIos, "")
	flag.StringVar(&c.Notification.AmazonSNS.ApplicationArnAndroid, "notification.amazonsns.applicationArnAndroid", c.Notification.AmazonSNS.ApplicationArnAndroid, "")
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

func (c *config) after() {
	if c.Datastore.Provider == "sqlite" {
		if c.Datastore.SQLite.Path == "" {
			c.Datastore.SQLite.Path = "/tmp/swagchat.db"
		}
		c.Datastore.Database = c.Datastore.SQLite.Path
	}

	if c.Datastore.Provider == "mysql" {
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxIdleConnection = "10"
		}
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxOpenConnection = "10"
		}
	}
}
