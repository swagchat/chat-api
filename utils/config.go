package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"strings"

	yaml "gopkg.in/yaml.v2"
)

type ctxKey int

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
	// HeaderUsername is http header for username
	HeaderUsername = "X-Preferred-Username"
	// HeaderRealm is http header for realm
	HeaderRealm = "X-Realm"
	// HeaderRealmRoles is http header for realm roles
	HeaderRealmRoles = "X-Realm-Roles"
	// HeaderAccountRoles is http header for account roles
	HeaderAccountRoles = "X-Account-Roles"

	CtxDsCfg ctxKey = iota
	CtxIsAppClient
	CtxUserID
	CtxRealm
	CtxRoomUser
	CtxSubscription
)

var (
	cfg           = NewConfig()
	IsShowVersion bool
)

type config struct {
	Version      string
	HTTPPort     string `yaml:"httpPort"`
	GRPCPort     string `yaml:"gRPCPort"`
	Profiling    bool
	DemoPage     bool `yaml:"demoPage"`
	ErrorLogging bool `yaml:"errorLogging"`
	Logging      *Logging
	PBroker      *PBroker
	SBroker      *SBroker
	Storage      *Storage
	Datastore    *Datastore
	Notification *Notification
	IdP          *IdP
}

type Logging struct {
	Level string
}

type PBroker struct {
	Provider string

	Direct struct {
		Endpoint string
	}

	Kafka struct {
		Host    string
		Port    string
		GroupID string `yaml:"groupId"`
		Topic   string
	}

	NSQ struct {
		Port           string
		NsqlookupdHost string
		NsqlookupdPort string
		NsqdHost       string
		NsqdPort       string
		Topic          string
		Channel        string
	}
}

type SBroker struct {
	Provider string

	Kafka struct {
		Host    string
		Port    string
		GroupID string `yaml:"groupId"`
		Topic   string
	}

	NSQ struct {
		Port           string
		NsqlookupdHost string
		NsqlookupdPort string
		NsqdHost       string
		NsqdPort       string
		Topic          string
		Channel        string
	}
}

type Storage struct {
	Provider string

	Local struct {
		Path string
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
		DirPath string `yaml:"dirPath"`
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

type IdP struct {
	Provider string

	// Keycloak
	Keycloak struct {
		BaseEndpoint string `yaml:"baseEndpoint"`
	}
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	c := &config{
		Version:      "0",
		HTTPPort:     "8101",
		GRPCPort:     "",
		Profiling:    false,
		DemoPage:     false,
		ErrorLogging: false,
		Logging:      &Logging{},
		PBroker:      &PBroker{},
		SBroker:      &SBroker{},
		Storage:      &Storage{},
		Datastore: &Datastore{
			Dynamic: false,
		},
		Notification: &Notification{},
		IdP:          &IdP{},
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
		c.HTTPPort = v
	}
	if v = os.Getenv("SWAG_HTTP_PORT"); v != "" {
		c.HTTPPort = v
	}
	if v = os.Getenv("SWAG_GRPC_PORT"); v != "" {
		c.GRPCPort = v
	}

	if v = os.Getenv("SWAG_PROFILING"); v != "" {
		if v == "true" {
			c.Profiling = true
		} else if v == "false" {
			c.Profiling = false
		}
	}
	if v = os.Getenv("SWAG_DEMO_PAGE"); v != "" {
		if v == "true" {
			c.DemoPage = true
		} else if v == "false" {
			c.DemoPage = false
		}
	}
	if v = os.Getenv("SWAG_ERROR_LOGGING"); v != "" {
		if v == "true" {
			c.ErrorLogging = true
		} else if v == "false" {
			c.ErrorLogging = false
		}
	}

	// Logging
	if v = os.Getenv("SWAG_LOGGING_LEVEL"); v != "" {
		c.Logging.Level = v
	}

	// PBroker
	if v = os.Getenv("SWAG_PBROKER_PROVIDER"); v != "" {
		c.PBroker.Provider = v
	}

	// PBroker - Direct
	if v = os.Getenv("SWAG_PBROKER_DIRECT_ENDPOINT"); v != "" {
		c.PBroker.Direct.Endpoint = v
	}

	// PBroker - Kafka
	if v = os.Getenv("SWAG_PBROKER_KAFKA_HOST"); v != "" {
		c.PBroker.Kafka.Host = v
	}
	if v = os.Getenv("SWAG_PBROKER_KAFKA_PORT"); v != "" {
		c.PBroker.Kafka.Port = v
	}
	if v = os.Getenv("SWAG_PBROKER_KAFKA_GROUPID"); v != "" {
		c.PBroker.Kafka.GroupID = v
	}
	if v = os.Getenv("SWAG_PBROKER_KAFKA_TOPIC"); v != "" {
		c.PBroker.Kafka.Topic = v
	}

	// PBroker - NSQ
	if v = os.Getenv("SWAG_PBROKER_NSQ_PORT"); v != "" {
		c.PBroker.NSQ.Port = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_NSQLOOKUPDHOST"); v != "" {
		c.PBroker.NSQ.NsqlookupdHost = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_NSQLOOKUPDPORT"); v != "" {
		c.PBroker.NSQ.NsqlookupdPort = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_NSQDHOST"); v != "" {
		c.PBroker.NSQ.NsqdHost = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_NSQDPORT"); v != "" {
		c.PBroker.NSQ.NsqdPort = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_TOPIC"); v != "" {
		c.PBroker.NSQ.Topic = v
	}
	if v = os.Getenv("SWAG_PBROKER_NSQ_CHANNEL"); v != "" {
		c.PBroker.NSQ.Channel = v
	}

	// SBroker
	if v = os.Getenv("SWAG_SBROKER_PROVIDER"); v != "" {
		c.SBroker.Provider = v
	}

	// SBroker - Kafka
	if v = os.Getenv("SWAG_SBROKER_KAFKA_HOST"); v != "" {
		c.SBroker.Kafka.Host = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_PORT"); v != "" {
		c.SBroker.Kafka.Port = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_GROUPID"); v != "" {
		c.SBroker.Kafka.GroupID = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_TOPIC"); v != "" {
		c.SBroker.Kafka.Topic = v
	}

	// SBroker - NSQ
	if v = os.Getenv("SWAG_SBROKER_NSQ_PORT"); v != "" {
		c.SBroker.NSQ.Port = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQLOOKUPDHOST"); v != "" {
		c.SBroker.NSQ.NsqlookupdHost = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQLOOKUPDPORT"); v != "" {
		c.SBroker.NSQ.NsqlookupdPort = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQDHOST"); v != "" {
		c.SBroker.NSQ.NsqdHost = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQDPORT"); v != "" {
		c.SBroker.NSQ.NsqdPort = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_TOPIC"); v != "" {
		c.SBroker.NSQ.Topic = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_CHANNEL"); v != "" {
		c.SBroker.NSQ.Channel = v
	}

	// Storage
	if v = os.Getenv("SWAG_STORAGE_PROVIDER"); v != "" {
		c.Storage.Provider = v
	}

	// Storage - Local
	if v = os.Getenv("SWAG_STORAGE_LOCAL_PATH"); v != "" {
		c.Storage.Local.Path = v
	}

	// Storage - Google Cloud Storage
	if v = os.Getenv("SWAG_STORAGE_GCS_PROJECT_ID"); v != "" {
		c.Storage.GCS.ProjectID = v
	}
	if v = os.Getenv("SWAG_STORAGE_GCS_JWT_PATH"); v != "" {
		c.Storage.GCS.JwtPath = v
	}
	if v = os.Getenv("SWAG_STORAGE_GCS_UPLOAD_BUCKET"); v != "" {
		c.Storage.GCS.UploadBucket = v
	}
	if v = os.Getenv("SWAG_STORAGE_GCS_UPLOAD_DIRECTORY"); v != "" {
		c.Storage.GCS.UploadDirectory = v
	}
	if v = os.Getenv("SWAG_STORAGE_GCS_THUMBNAIL_BUCKET"); v != "" {
		c.Storage.GCS.ThumbnailBucket = v
	}
	if v = os.Getenv("SWAG_STORAGE_GCS_THUMBNAIL_DIRECTORY"); v != "" {
		c.Storage.GCS.ThumbnailDirectory = v
	}

	// Storage - AWS S3
	if v = os.Getenv("SWAG_STORAGE_AWSS3_REGION"); v != "" {
		c.Storage.AWSS3.Region = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_ACCESS_KEY_ID"); v != "" {
		c.Storage.AWSS3.AccessKeyID = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_SECRET_ACCESS_KEY"); v != "" {
		c.Storage.AWSS3.SecretAccessKey = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_UPLOAD_BUCKET"); v != "" {
		c.Storage.AWSS3.UploadBucket = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_UPLOAD_DIRECTORY"); v != "" {
		c.Storage.AWSS3.UploadDirectory = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_THUMBNAIL_BUCKET"); v != "" {
		c.Storage.AWSS3.ThumbnailBucket = v
	}
	if v = os.Getenv("SWAG_STORAGE_AWSS3_THUMBNAIL_DIRECTORY"); v != "" {
		c.Storage.AWSS3.ThumbnailDirectory = v
	}

	// Datastore
	if v = os.Getenv("SWAG_DATASTORE_DYNAMIC"); v != "" {
		if v == "true" {
			c.Datastore.Dynamic = true
		} else if v == "false" {
			c.Datastore.Dynamic = false
		}
	}
	if v = os.Getenv("SWAG_DATASTORE_PROVIDER"); v != "" {
		c.Datastore.Provider = v
	}

	if v = os.Getenv("SWAG_DATASTORE_USER"); v != "" {
		c.Datastore.User = v
	}
	if v = os.Getenv("SWAG_DATASTORE_PASSWORD"); v != "" {
		c.Datastore.Password = v
	}
	if v = os.Getenv("SWAG_DATASTORE_DATABASE"); v != "" {
		c.Datastore.Database = v
	}
	if v = os.Getenv("SWAG_DATASTORE_TABLE_NAME_PREFIX"); v != "" {
		c.Datastore.TableNamePrefix = v
	}
	if v = os.Getenv("SWAG_DATASTORE_MAX_IDLE_CONNECTION"); v != "" {
		c.Datastore.MaxIdleConnection = v
	}
	if v = os.Getenv("SWAG_DATASTORE_MAX_OPEN_CONNECTION"); v != "" {
		c.Datastore.MaxOpenConnection = v
	}

	var master *ServerInfo
	mHost := os.Getenv("SWAG_DATASTORE_MASTER_HOST")
	mPort := os.Getenv("SWAG_DATASTORE_MASTER_PORT")
	if mHost != "" && mPort != "" {
		master = &ServerInfo{}
		master.Host = mHost
		master.Port = mPort
		c.Datastore.Master = master
		mServerName := os.Getenv("SWAG_DATASTORE_MASTER_SERVER_NAME")
		mServerCaPath := os.Getenv("SWAG_DATASTORE_MASTER_SERVER_CA_PATH")
		mClientCertPath := os.Getenv("SWAG_DATASTORE_MASTER_CLIENT_CERT_PATH")
		mClientKeyPath := os.Getenv("SWAG_DATASTORE_MASTER_CLIENT_KEY_PATH")
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
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_HOSTS"); v != "" {
		rHosts = strings.Split(v, ",")
	}
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_PORTS"); v != "" {
		rPorts = strings.Split(v, ",")
	}
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_SERVER_NAMES"); v != "" {
		rServerName = strings.Split(v, ",")
	}
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_SERVER_CA_PATHS"); v != "" {
		rServerCaPath = strings.Split(v, ",")
	}
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_CLIENT_CERT_PATHS"); v != "" {
		rClientCertPath = strings.Split(v, ",")
	}
	if v = os.Getenv("SWAG_DATASTORE_REPLICA_CLIENT_KEY_PATHS"); v != "" {
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
	if v = os.Getenv("SWAG_DATASTORE_SQLITE_DIRPATH"); v != "" {
		c.Datastore.SQLite.DirPath = v
	}

	// Notification
	if v = os.Getenv("SWAG_NOTIFICATION_PROVIDER"); v != "" {
		c.Notification.Provider = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_ROOM_TOPIC_NAME_PREFIX"); v != "" {
		c.Notification.RoomTopicNamePrefix = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_DEFAULT_BADGE_COUNT"); v != "" {
		c.Notification.DefaultBadgeCount = v
	}

	// Notification - Amazon SNS
	if v = os.Getenv("SWAG_NOTIFICATION_AMAZONSNS_REGION"); v != "" {
		c.Notification.AmazonSNS.Region = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_AMAZONSNS_ACCESS_KEY_ID"); v != "" {
		c.Notification.AmazonSNS.AccessKeyID = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_AMAZONSNS_SECRET_ACCESS_KEY"); v != "" {
		c.Notification.AmazonSNS.SecretAccessKey = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_AMAZONSNS_APPLICATION_ARN_IOS"); v != "" {
		c.Notification.AmazonSNS.ApplicationArnIos = v
	}
	if v = os.Getenv("SWAG_NOTIFICATION_AMAZONSNS_APPLICATION_ARN_ANDROID"); v != "" {
		c.Notification.AmazonSNS.ApplicationArnAndroid = v
	}

	// IdP
	if v = os.Getenv("SWAG_IDP_PROVIDER"); v != "" {
		c.IdP.Provider = v
	}

	// IdP Keycloak
	if v = os.Getenv("SWAG_IDP_KEYCLOAK_BASEENDPOINT"); v != "" {
		c.IdP.Keycloak.BaseEndpoint = v
	}
}

func (c *config) parseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HTTPPort, "httpPort", c.HTTPPort, "")
	flag.StringVar(&c.GRPCPort, "grpcPort", c.GRPCPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	// PBroker
	flag.StringVar(&c.PBroker.Provider, "pbroker.provider", c.PBroker.Provider, "")

	// PBroker - Direct
	flag.StringVar(&c.PBroker.Direct.Endpoint, "pbroker.direct.endpoint", c.PBroker.Direct.Endpoint, "")

	// PBroker - kafka
	flag.StringVar(&c.PBroker.Kafka.Host, "pbroker.kafka.host", c.PBroker.Kafka.Host, "")
	flag.StringVar(&c.PBroker.Kafka.Port, "pbroker.kafka.port", c.PBroker.Kafka.Port, "")
	flag.StringVar(&c.PBroker.Kafka.GroupID, "pbroker.kafka.groupId", c.PBroker.Kafka.GroupID, "")
	flag.StringVar(&c.PBroker.Kafka.Topic, "pbroker.kafka.topic", c.PBroker.Kafka.Topic, "")

	// PBroker - NSQ
	flag.StringVar(&c.PBroker.NSQ.NsqlookupdHost, "pbroker.nsq.nsqlookupdHost", c.PBroker.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.PBroker.NSQ.NsqlookupdPort, "pbroker.nsq.nsqlookupdPort", c.PBroker.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.PBroker.NSQ.NsqdHost, "pbroker.nsq.nsqdHost", c.PBroker.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.PBroker.NSQ.NsqdPort, "pbroker.nsq.nsqdPort", c.PBroker.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.PBroker.NSQ.Topic, "pbroker.nsq.topic", c.PBroker.NSQ.Topic, "Topic name")
	flag.StringVar(&c.PBroker.NSQ.Channel, "pbroker.nsq.channel", c.PBroker.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

	// SBroker
	flag.StringVar(&c.SBroker.Provider, "sbroker.provider", c.SBroker.Provider, "")

	// SBroker - kafka
	flag.StringVar(&c.SBroker.Kafka.Host, "sbroker.kafka.host", c.SBroker.Kafka.Host, "")
	flag.StringVar(&c.SBroker.Kafka.Port, "sbroker.kafka.port", c.SBroker.Kafka.Port, "")
	flag.StringVar(&c.SBroker.Kafka.GroupID, "sbroker.kafka.groupId", c.SBroker.Kafka.GroupID, "")
	flag.StringVar(&c.SBroker.Kafka.Topic, "sbroker.kafka.topic", c.SBroker.Kafka.Topic, "")

	// SBroker - NSQ
	flag.StringVar(&c.SBroker.NSQ.NsqlookupdHost, "sbroker.nsq.nsqlookupdHost", c.SBroker.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.SBroker.NSQ.NsqlookupdPort, "sbroker.nsq.nsqlookupdPort", c.SBroker.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.SBroker.NSQ.NsqdHost, "sbroker.nsq.nsqdHost", c.SBroker.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.SBroker.NSQ.NsqdPort, "sbroker.nsq.nsqdPort", c.SBroker.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.SBroker.NSQ.Topic, "sbroker.nsq.topic", c.SBroker.NSQ.Topic, "Topic name")
	flag.StringVar(&c.SBroker.NSQ.Channel, "sbroker.nsq.channel", c.SBroker.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

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
	flag.StringVar(&c.Datastore.SQLite.DirPath, "datastore.sqlite.dirPath", c.Datastore.SQLite.DirPath, "")

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

	// IdP
	flag.StringVar(&c.IdP.Provider, "idp.provider", c.IdP.Provider, "")

	// IdP Keycloak
	flag.StringVar(&c.IdP.Keycloak.BaseEndpoint, "idp.keycloak.baseEndpoint", c.IdP.Keycloak.BaseEndpoint, "")
}

func (c *config) after() {
	if c.Logging.Level == "" {
		c.Logging.Level = "development"
	}

	if c.Storage.Provider == "" {
		c.Storage.Provider = "local"
	}
	if c.Storage.Provider == "local" && c.Storage.Local.Path == "" {
		c.Storage.Local.Path = "data/assets"
	}

	if c.Datastore.Provider == "" {
		c.Datastore.Provider = "sqlite"
	}
	if c.Datastore.Provider == "sqlite" {
		if c.Datastore.SQLite.DirPath == "" {
			c.Datastore.SQLite.DirPath = "/tmp"
		}
	}

	if c.Datastore.Provider == "mysql" {
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxIdleConnection = "10"
		}
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxOpenConnection = "10"
		}
	}

	if c.Datastore.Database == "" {
		c.Datastore.Database = "swagchat"
	}

	if c.IdP.Provider == "" {
		c.IdP.Provider = "local"
	}
}
