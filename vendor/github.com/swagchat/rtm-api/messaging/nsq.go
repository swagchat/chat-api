package messaging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"unsafe"

	nsq "github.com/nsqio/go-nsq"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

var Con *nsq.Consumer

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
func hello(s string) {
	fmt.Println(s)
}

func Subscribe() {
	if utils.Que.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := utils.Que.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}
		Con, _ = nsq.NewConsumer(utils.Que.Topic, channel, config)
		Con.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			log.Printf("[NSQ]Got a message: %v", message)
			services.Srv.Broadcast <- message.Body
			return nil
		}))
		err = Con.ConnectToNSQLookupd(utils.Que.NsqlookupdHost + ":" + utils.Que.NsqlookupdPort)
		if err != nil {
			log.Panic("Could not connect")
		}
	}
}

func Unsubscribe() {
	if Con != nil {
		hostname, err := os.Hostname()
		resp, err := http.Post("http://"+utils.Que.NsqdHost+":"+utils.Que.NsqdPort+"/channel/delete?topic="+utils.Que.Topic+"&channel="+hostname, "text/plain", nil)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp)
	}
}
