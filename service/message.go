package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/producer"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

// SendMessage creates message
func SendMessage(ctx context.Context, req *model.SendMessageRequest) (*model.Message, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("SendMessage", "service")
	defer tracer.Provider(ctx).Finish(span)

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	room, errRes := confirmRoomExist(ctx, *req.RoomID)
	if errRes != nil {
		errRes.Message = "Failed to create message."
		return nil, errRes
	}

	user, errRes := confirmUserExist(ctx, *req.UserID, datastore.SelectUserOptionWithRoles(true))
	if errRes != nil {
		errRes.Message = "Failed to create message."
		return nil, errRes
	}

	message := req.GenerateMessage()

	if message.Type == model.MessageTypeIndicatorStart || message.Type == model.MessageTypeIndicatorEnd {
		publishMessage(ctx, message)
		return nil, nil
	}

	_, errRes = confirmMessageNotExist(ctx, message.MessageID)
	if errRes != nil {
		errRes.Message = "Failed to create message."
		return nil, errRes
	}

	err := datastore.Provider(ctx).InsertMessage(message)
	if err != nil {
		errRes := model.NewErrorResponse("Failed to create message.", http.StatusInternalServerError, model.WithError(err))
		return nil, errRes
	}

	// notification
	lastMessage := "" // TODO
	mi := &notification.MessageInfo{
		Text: fmt.Sprintf("[%s]%s", room.Name, lastMessage),
	}
	cfg := config.Config()
	if cfg.Notification.DefaultBadgeCount != "" {
		dBadgeCount, err := strconv.Atoi(cfg.Notification.DefaultBadgeCount)
		if err == nil {
			mi.Badge = dBadgeCount
		}
	}
	go notification.Provider(ctx).Publish(room.NotificationTopicID, room.RoomID, mi)

	publishMessage(ctx, message)
	webhookMessage(ctx, message, user)

	return message, nil
}

// RetrieveMessage gets message
func RetrieveMessage(ctx context.Context, messageID string) (*model.Message, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("RetrieveMessage", "service")
	defer tracer.Provider(ctx).Finish(span)

	if messageID == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: "messageId is required, but it's empty.",
			},
		}
		return nil, model.NewErrorResponse("Failed to retrieve message.", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve message.", http.StatusInternalServerError, model.WithError(err))
	}
	if message == nil {
		return nil, model.NewErrorResponse("", http.StatusNotFound)
	}

	return message, nil
}

func publishMessage(ctx context.Context, message *model.Message) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
		datastore.SelectUserIDsOfRoomUserOptionWithRoomID(message.RoomID),
		datastore.SelectUserIDsOfRoomUserOptionWithRoles([]int32{message.Role}),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(message)
	event := &scpb.EventData{
		Type:    scpb.EventType_MessageEvent,
		Data:    buffer.Bytes(),
		UserIDs: userIDs,
	}
	err = producer.Provider(ctx).PublishMessage(event)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
