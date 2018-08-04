package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/model"
)

func TestMessageStore(t *testing.T) {
	newMessage := &model.Message{}
	newMessage.MessageID = "message-id"
	newMessage.RoomID = "datastore-room-id-0001"
	newMessage.UserID = "datastore-user-id-0001"
	newMessage.Type = "text"
	newMessage.Payload = []byte(`{"text":"test"}`)
	newMessage.Role = config.RoleGeneral
	newMessage.Created = 123456789
	newMessage.Modified = 123456789
	err := Provider(ctx).InsertMessage(newMessage)
	if err != nil {
		t.Fatalf("failed insert message test")
	}

	messages, err := Provider(ctx).SelectMessages(10, 0)
	if err != nil {
		t.Fatalf("failed select messages test")
	}
	if len(messages) != 1 {
		t.Fatalf("failed select messages test")
	}

	message, err := Provider(ctx).SelectMessage("message-id")
	if err != nil {
		t.Fatalf("failed select message test")
	}
	if message == nil {
		t.Fatalf("failed select message test")
	}

	message.Payload = []byte(`{"text":"test-update"}`)
	err = Provider(ctx).UpdateMessage(message)
	if err != nil {
		t.Fatalf("failed update message test")
	}

	updatedMessage, err := Provider(ctx).SelectMessage("message-id")
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
