package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type ctxKey int

const (
	// AppName is Application name
	AppName = "operator-api"
	// APIVersion is API version
	APIVersion = "0"
	// BuildVersion is API build version
	BuildVersion = "0.1.0"

	// HeaderUserID is http header for userID
	HeaderUserID = "X-Sub"
	// HeaderUsername is http header for username
	HeaderUsername = "X-Preferred-Username"
	// HeaderWorkspace is http header for workspace
	HeaderWorkspace = "X-Realm"
	// HeaderRealmRoles is http header for roles
	HeaderRealmRoles = "X-Realm-Roles"
	// HeaderAccountRoles is http header for account roles
	HeaderAccountRoles = "X-Account-Roles"

	CtxUserID ctxKey = iota
	CtxWorkspace

	RoleGeneral  int32 = 1
	RoleGuest    int32 = 2
	RoleOperator int32 = 3
	RoleBot      int32 = 5
)

var (
	cfg         = NewConfig()
	showVersion = false
	showHelp    = false
	// StopRun is a flag for stop run server
	StopRun = false
)

type config struct {
	Version                        string
	HTTPPort                       string `yaml:"httpPort"`
	GRPCPort                       string `yaml:"grpcPort"`
	ChatAPIGRPCEndpoint            string `yaml:"chatApiGrpcEndpoint"`
	MessageConnectorAPIRPCEndpoint string `yaml:"messageConnectorApiGrpcEndpoint"`
	Profiling                      bool
	EnableDeveloperMessage         bool `yaml:"enableDeveloperMessage"`
	Logger                         *Logger
	Message                        *Message
	PBroker                        *PBroker
	Datastore                      *Datastore
	IdP                            *IdP
}

// Logger is settings of logger
type Logger struct {
	// EnableConsole is a flag for enable console log.
	EnableConsole bool `yaml:"enableConsole"`
	// ConsoleFormat is a format for console log.
	ConsoleFormat string `yaml:"consoleFormat"`
	// ConsoleLevel is a level for console log.
	ConsoleLevel string `yaml:"consoleLevel"`
	// EnableFile is a flag for enable file log.
	EnableFile bool `yaml:"enableFile"`
	// FileFormat is a format for file log.
	FileFormat string `yaml:"fileFormat"`
	// FileLevel is a log level for file log.
	FileLevel string `yaml:"fileLevel"`
	// FilePath is a file path for file log.
	FilePath string `yaml:"filePath"`
}

type Message struct {
	Provider string // chat-api, message-connector-api
	Protocol string // http, grpc, gossip
	Endpoint string
	Topic    string
}

type PBroker struct {
	Provider string

	MessageConnector struct {
		Endpoint string
		Protocol string // http, grpc
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
	EnableLogging     bool `yaml:"enableLogging"`
	SQLite            *SQLite
}

type SQLite struct {
	OnMemory bool   `yaml:"onMemory"`
	DirPath  string `yaml:"dirPath"`
}

type ServerInfo struct {
	Host           string
	Port           string
	ServerName     string `yaml:"serverName"`
	ServerCaPath   string `yaml:"serverCaPath"`
	ClientCertPath string `yaml:"clientCertPath"`
	ClientKeyPath  string `yaml:"clientKeyPath"`
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

	c := defaultSetting()
	c.loadEnv()

	err := c.parseFlag(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to load setting. %v", err)
	}

	err = c.validate()
	if err != nil {
		log.Fatalf("Invalid setting. %v", err)
	}

	err = c.after()
	if err != nil {
		log.Fatalf("%v", err)
	}
	return c
}

func Config() *config {
	return cfg
}

func defaultSetting() *config {
	sqlite := &SQLite{
		OnMemory: false,
		DirPath:  "",
	}
	return &config{
		HTTPPort:               "8106",
		GRPCPort:               "9106",
		Profiling:              false,
		EnableDeveloperMessage: false,
		Logger: &Logger{
			EnableConsole: true,
			ConsoleFormat: "text",
			ConsoleLevel:  "debug",
			EnableFile:    false,
		},
		Message: &Message{},
		PBroker: &PBroker{},
		Datastore: &Datastore{
			Dynamic:       false,
			Provider:      "sqlite",
			Database:      "operator",
			EnableLogging: false,
			SQLite:        sqlite,
		},
		IdP: &IdP{},
	}
}

func (c *config) loadYaml(buf []byte) {
	yaml.Unmarshal(buf, c)
}

func (c *config) loadEnv() {
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
	if v = os.Getenv("SWAG_CHATAPIGRPC_ENDPOINT"); v != "" {
		c.ChatAPIGRPCEndpoint = v
	}
	if v = os.Getenv("SWAG_MESSAGECONNECTORAPIGRPC_ENDPOINT"); v != "" {
		c.MessageConnectorAPIRPCEndpoint = v
	}
	if v = os.Getenv("SWAG_PROFILING"); v == "true" {
		c.Profiling = true
	}
	if v = os.Getenv("SWAG_ENABLE_DEVELOPER_MESSAGE"); v == "true" {
		c.EnableDeveloperMessage = true
	}

	// Logger
	if v = os.Getenv("SWAG_LOGGER_ENABLE_CONSOLE"); v == "true" {
		c.Logger.EnableConsole = true
	}
	if v = os.Getenv("SWAG_LOGGER_CONSOLE_FORMAT"); v != "" {
		c.Logger.ConsoleFormat = v
	}
	if v = os.Getenv("SWAG_LOGGER_CONSOLE_LEVEL"); v != "" {
		c.Logger.ConsoleLevel = v
	}
	if v = os.Getenv("SWAG_LOGGER_ENABLE_FILE"); v == "true" {
		c.Logger.EnableFile = true
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_FORMAT"); v != "" {
		c.Logger.FileFormat = v
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_LEVEL"); v != "" {
		c.Logger.FileLevel = v
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_PATH"); v != "" {
		c.Logger.FilePath = v
	}

	// Message
	if v = os.Getenv("SWAG_MESSAGE_PROVIDER"); v != "" {
		c.Message.Provider = v
	}
	if v = os.Getenv("SWAG_MESSAGE_PROTOCOL"); v != "" {
		c.Message.Protocol = v
	}
	if v = os.Getenv("SWAG_MESSAGE_ENDPOINT"); v != "" {
		c.Message.Endpoint = v
	}

	// PBroker
	if v = os.Getenv("SWAG_PBROKER_PROVIDER"); v != "" {
		c.PBroker.Provider = v
	}

	// PBroker - MessageConnector
	if v = os.Getenv("SWAG_PBROKER_MESSAGECONNECTOR_ENDPOINT"); v != "" {
		c.PBroker.MessageConnector.Endpoint = v
	}
	if v = os.Getenv("SWAG_PBROKER_MESSAGECONNECTOR_PROTOCOL"); v != "" {
		c.PBroker.MessageConnector.Protocol = v
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
	if v = os.Getenv("SWAG_DATASTORE_SQLITE_ONMEMORY"); v == "true" {
		c.Datastore.SQLite.OnMemory = true
	}
	if v = os.Getenv("SWAG_DATASTORE_SQLITE_DIRPATH"); v != "" {
		c.Datastore.SQLite.DirPath = v
	}

	// IDP
	if v = os.Getenv("SWAG_IDP_PROVIDER"); v != "" {
		c.IdP.Provider = v
	}
	if v = os.Getenv("SWAG_IDP_KEYCLOAK_BASEENDPOINT"); v != "" {
		c.IdP.Keycloak.BaseEndpoint = v
	}
}

func (c *config) parseFlag(args []string) error {
	if len(args) == 0 {
		return nil
	}

	flags := flag.NewFlagSet("ChatAPI Flags", flag.ContinueOnError)

	flags.BoolVar(&showVersion, "v", false, "show version")
	flags.BoolVar(&showVersion, "version", false, "show version")
	flags.BoolVar(&showHelp, "h", false, "show help")
	flags.BoolVar(&showHelp, "help", false, "show help")

	flags.StringVar(&c.HTTPPort, "httpPort", c.HTTPPort, "")
	flags.StringVar(&c.GRPCPort, "grpcPort", c.GRPCPort, "")

	flags.StringVar(&c.ChatAPIGRPCEndpoint, "chatApiGrpcEndpoint", c.ChatAPIGRPCEndpoint, "")
	flags.StringVar(&c.MessageConnectorAPIRPCEndpoint, "messageConnectorApiGrpcEndpoint", c.MessageConnectorAPIRPCEndpoint, "")

	var profiling string
	flags.StringVar(&profiling, "profiling", "", "")

	var enableDeveloperMessage string
	flags.StringVar(&enableDeveloperMessage, "enableDeveloperMessage", "", "false")

	// Logging
	flags.BoolVar(&c.Logger.EnableConsole, "logger.enableConsole", c.Logger.EnableConsole, "")
	flags.StringVar(&c.Logger.ConsoleFormat, "logger.consoleFormat", c.Logger.ConsoleFormat, "")
	flags.StringVar(&c.Logger.ConsoleLevel, "logger.consoleLevel", c.Logger.ConsoleLevel, "")
	flags.BoolVar(&c.Logger.EnableFile, "logger.enableFile", c.Logger.EnableFile, "")
	flags.StringVar(&c.Logger.FileFormat, "logger.fileFormat", c.Logger.FileFormat, "")
	flags.StringVar(&c.Logger.FileLevel, "logger.fileLevel", c.Logger.FileLevel, "")
	flags.StringVar(&c.Logger.FilePath, "logger.filePath", c.Logger.FilePath, "")

	// Message
	flags.StringVar(&c.Message.Provider, "message.provider", c.Message.Provider, "")
	flags.StringVar(&c.Message.Protocol, "message.protocol", c.Message.Protocol, "")
	flags.StringVar(&c.Message.Endpoint, "message.endpoint", c.Message.Endpoint, "")

	// PBroker
	flags.StringVar(&c.PBroker.Provider, "pbroker.provider", c.PBroker.Provider, "")

	// PBroker - MessageConnector
	flags.StringVar(&c.PBroker.MessageConnector.Endpoint, "pbroker.messageConnector.endpoint", c.PBroker.MessageConnector.Endpoint, "")
	flags.StringVar(&c.PBroker.MessageConnector.Protocol, "pbroker.messageConnector.protocol", c.PBroker.MessageConnector.Protocol, "")

	// PBroker - kafka
	flags.StringVar(&c.PBroker.Kafka.Host, "pbroker.kafka.host", c.PBroker.Kafka.Host, "")
	flags.StringVar(&c.PBroker.Kafka.Port, "pbroker.kafka.port", c.PBroker.Kafka.Port, "")
	flags.StringVar(&c.PBroker.Kafka.GroupID, "pbroker.kafka.groupId", c.PBroker.Kafka.GroupID, "")
	flags.StringVar(&c.PBroker.Kafka.Topic, "pbroker.kafka.topic", c.PBroker.Kafka.Topic, "")

	// PBroker - NSQ
	flags.StringVar(&c.PBroker.NSQ.NsqlookupdHost, "pbroker.nsq.nsqlookupdHost", c.PBroker.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flags.StringVar(&c.PBroker.NSQ.NsqlookupdPort, "pbroker.nsq.nsqlookupdPort", c.PBroker.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flags.StringVar(&c.PBroker.NSQ.NsqdHost, "pbroker.nsq.nsqdHost", c.PBroker.NSQ.NsqdHost, "Host name of nsqd")
	flags.StringVar(&c.PBroker.NSQ.NsqdPort, "pbroker.nsq.nsqdPort", c.PBroker.NSQ.NsqdPort, "Port no of nsqd")
	flags.StringVar(&c.PBroker.NSQ.Topic, "pbroker.nsq.topic", c.PBroker.NSQ.Topic, "Topic name")
	flags.StringVar(&c.PBroker.NSQ.Channel, "pbroker.nsq.channel", c.PBroker.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

	// Datastore
	flags.BoolVar(&c.Datastore.Dynamic, "datastore.dynamic", c.Datastore.Dynamic, "")
	flags.StringVar(&c.Datastore.Provider, "datastore.provider", c.Datastore.Provider, "")
	flags.StringVar(&c.Datastore.TableNamePrefix, "datastore.tableNamePrefix", c.Datastore.TableNamePrefix, "")
	flags.StringVar(&c.Datastore.User, "datastore.user", c.Datastore.User, "")
	flags.StringVar(&c.Datastore.Password, "datastore.password", c.Datastore.Password, "")
	flags.StringVar(&c.Datastore.Database, "datastore.database", c.Datastore.Database, "")
	flags.StringVar(&c.Datastore.MaxIdleConnection, "datastore.maxIdleConnection", c.Datastore.MaxIdleConnection, "")
	flags.StringVar(&c.Datastore.MaxOpenConnection, "datastore.maxOpenConnection", c.Datastore.MaxOpenConnection, "")

	var (
		mHostStr           string
		mPortsStr          string
		mServerNameStr     string
		mServerCaPathStr   string
		mClientCertPathStr string
		mClientKeyPathStr  string
	)
	flags.StringVar(&mHostStr, "datastore.masterHost", mHostStr, "")
	flags.StringVar(&mPortsStr, "datastore.masterPort", mPortsStr, "")
	if mHostStr != "" && mPortsStr != "" {
		c.Datastore.Master = &ServerInfo{
			Host: mHostStr,
			Port: mPortsStr,
		}
		flags.StringVar(&mServerNameStr, "datastore.masterServerName", mServerNameStr, "")
		flags.StringVar(&mServerCaPathStr, "datastore.masterServerCaPath", mServerCaPathStr, "")
		flags.StringVar(&mClientCertPathStr, "datastore.masterClientCertPath", mClientCertPathStr, "")
		flags.StringVar(&mClientKeyPathStr, "datastore.masterClientKeyPath", mClientKeyPathStr, "")
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
	flags.StringVar(&rHostsStr, "datastore.replicaHosts", rHostsStr, "")
	flags.StringVar(&rPortsStr, "datastore.replicaPorts", rPortsStr, "")
	flags.StringVar(&rServerNamesStr, "datastore.replicaServerNames", rServerNamesStr, "")
	flags.StringVar(&rServerCaPathsStr, "datastore.replicaServerCaPaths", rServerCaPathsStr, "")
	flags.StringVar(&rClientCertPathsStr, "datastore.replicaClientCertPaths", rClientCertPathsStr, "")
	flags.StringVar(&rClientKeyPathsStr, "datastore.replicaClientKeyPaths", rClientKeyPathsStr, "")

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
	var onMemory string
	flags.StringVar(&onMemory, "datastore.sqlite.onMemory", "", "false")
	flags.StringVar(&c.Datastore.SQLite.DirPath, "datastore.sqlite.dirPath", c.Datastore.SQLite.DirPath, "")

	// IdP
	flags.StringVar(&c.IdP.Provider, "idp.provider", c.IdP.Provider, "")

	// IdP Keycloak
	flags.StringVar(&c.IdP.Keycloak.BaseEndpoint, "idp.keycloak.baseEndpoint", c.IdP.Keycloak.BaseEndpoint, "")

	configPath := ""
	flags.StringVar(&configPath, "config", "", "config file(yaml format)")

	err := flags.Parse(args)
	if err != nil {
		if flag.Lookup("test.") != nil { // for testing
			return errors.Wrap(err, "")
		}
	}

	if showHelp {
		flags.PrintDefaults()
		StopRun = true
		return nil
	}

	if showVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", APIVersion, BuildVersion)
		StopRun = true
		return nil
	}

	if configPath != "" {
		if !isExists(configPath) {
			return fmt.Errorf("File not found [%s]", configPath)
		}
		buf, _ := ioutil.ReadFile(configPath)
		c.loadYaml(buf)
	}

	if profiling == "true" {
		c.Profiling = true
	}

	if enableDeveloperMessage == "true" {
		c.EnableDeveloperMessage = true
	}

	if onMemory == "true" {
		c.Datastore.SQLite.OnMemory = true
	}

	return nil
}

func isExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (c *config) validate() error {
	// Logger
	if c.Logger.EnableConsole {
		f := c.Logger.ConsoleFormat
		if f == "" || !(f == "text" || f == "json") {
			return errors.New("Please set logger.consoleFormat to \"text\" or \"json\"")
		}
		l := c.Logger.ConsoleLevel
		if l == "" || !(l == "debug" || l == "info" || l == "warn" || l == "error") {
			return errors.New("Please set logger.consoleLevel to \"debug\" or \"info\" or \"warn\" or \"error\"")
		}
	}
	if c.Logger.EnableFile {
		f := c.Logger.FileFormat
		if f == "" || !(f == "text" || f == "json") {
			return errors.New("Please set logger.fileFormat to \"text\" or \"json\"")
		}
		l := c.Logger.FileLevel
		if l == "" || !(l == "debug" || l == "info" || l == "warn" || l == "error") {
			return errors.New("Please set logger.fileLevel to \"debug\" or \"info\" or \"warn\" or \"error\"")
		}
		if c.Logger.FilePath == "" {
			return errors.New("Please set logger.filePath")
		}
	}

	// Idp
	if c.IdP.Provider == "keycloak" {
		if c.IdP.Keycloak.BaseEndpoint == "" {
			return errors.Wrap(errors.New("keycloak base endpoint is empty"), "")
		}
	}

	return nil
}

func (c *config) after() error {
	if c.Datastore.Provider == "" {
		c.Datastore.Provider = "sqlite"
	}

	if c.Datastore.Provider == "sqlite" && !c.Datastore.SQLite.OnMemory && c.Datastore.SQLite.DirPath == "" {
		tmpDirPath, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		c.Datastore.SQLite.DirPath = tmpDirPath
	}

	if c.Datastore.Provider == "mysql" {
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxIdleConnection = "10"
		}
		if c.Datastore.MaxIdleConnection == "" {
			c.Datastore.MaxOpenConnection = "10"
		}
	}

	return nil
}
