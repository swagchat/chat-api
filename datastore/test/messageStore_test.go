package datastore_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

func TestMessageStore(t *testing.T) {
	newMessage := &model.Message{}
	newMessage.MessageID = "message-id"
	newMessage.RoomID = "room-id-0000"
	newMessage.UserID = "user-id-0000"
	newMessage.Type = "text"
	newMessage.Payload = []byte(`{"text":"test"}`)
	newMessage.Role = utils.RoleGeneral
	newMessage.Created = 123456789
	newMessage.Modified = 123456789
	err := datastore.Provider(ctx).InsertMessage(newMessage)
	if err != nil {
		t.Fatalf("failed insert message test")
	}

	messages, err := datastore.Provider(ctx).SelectMessages(10, 0)
	if err != nil {
		t.Fatalf("failed select messages test")
	}
	if len(messages) != 1 {
		t.Fatalf("failed select messages test")
	}

	message, err := datastore.Provider(ctx).SelectMessage("message-id")
	if err != nil {
		t.Fatalf("failed select message test")
	}
	if message == nil {
		t.Fatalf("failed select message test")
	}

	message.Payload = []byte(`{"text":"test-update"}`)
	err = datastore.Provider(ctx).UpdateMessage(message)
	if err != nil {
		t.Fatalf("failed update message test")
	}

	updatedMessage, err := datastore.Provider(ctx).SelectMessage("message-id")
	if err != nil {
		t.Fatalf("failed select message test")
	}

	bytes, err := updatedMessage.Payload.MarshalJSON()
	if err != nil {
		t.Fatalf("failed marshal payload")
	}
	if string(bytes) != `{"text":"test-update"}` {
		t.Fatalf("failed update message test")
	}
}
