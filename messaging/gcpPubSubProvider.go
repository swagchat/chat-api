package messaging

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	pubsub "google.golang.org/api/pubsub/v1"
)

type GcpPubSubProvider struct {
	thumbnailTopic    string
	pushMessageTopic  string
	scope             string
	jwtConfigFilepath string
}

var gcpPubSubService *pubsub.Service

func (provider GcpPubSubProvider) Init() error {
	if gcpPubSubService == nil {
		data, err := ioutil.ReadFile(provider.jwtConfigFilepath)
		if err != nil {
			return err
		}

		conf, err := google.JWTConfigFromJSON(data, provider.scope)
		if err != nil {
			return err
		}
		client := conf.Client(oauth2.NoContext)

		service, err := pubsub.New(client)
		if err != nil {
			return err
		}
		gcpPubSubService = service
	}
	return nil
}

func (provider GcpPubSubProvider) PublishMessage(messagingInfo *MessagingInfo) error {
	// REF https://github.com/google/google-api-go-client/blob/master/examples/pubsub.go
	bytes, err := json.Marshal(messagingInfo.Message)
	if err != nil {
		log.Println(err.Error())
	}
	pubsubMessage := &pubsub.PubsubMessage{
		Data: base64.StdEncoding.EncodeToString(bytes),
	}
	publishRequest := &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{pubsubMessage},
	}
	if _, err := gcpPubSubService.Projects.Topics.Publish(provider.pushMessageTopic, publishRequest).Do(); err != nil {
		log.Fatalf("connectIRC Publish().Do() failed: %v", err)
	}
	log.Println("Published a message to the topic.")

	return nil
}
