package utils

const (
	API_VERSION = "v0"
)

var (
	Realtime RealtimeSetting
	Que      QueSetting
)

type RealtimeSetting struct {
	Port string
}

type QueSetting struct {
	Port           string
	NsqlookupdHost string
	NsqlookupdPort string
	NsqdHost       string
	NsqdPort       string
	Topic          string
	Channel        string
}
