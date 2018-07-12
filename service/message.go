package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/utils"
)

// PostMessage is post message
func PostMessage(ctx context.Context, posts *model.Messages) *model.ResponseMessages {
	messageIds := make([]string, 0)
	pds := make([]*model.ProblemDetail, 0)
	for _, post := range posts.Messages {
		logger.Info(fmt.Sprintf("Start CreateMessage. Message=[%#v]", post))
		room, pd := selectRoom(ctx, post.RoomID)
		if pd != nil {
			pds = append(pds, &model.ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "roomId",
						Reason: "roomId is invalid. Not exist room.",
					},
				},
				Status: http.StatusBadRequest,
			})
			continue
		}

		user, pd := selectUser(ctx, post.UserID, datastore.WithRoles(true))
		if pd != nil {
			pds = append(pds, &model.ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "userId",
						Reason: "userId is invalid. Not exist user.",
					},
				},
				Status: http.StatusBadRequest,
			})
			continue
		}

		if pd := post.IsValid(); pd != nil {
			pds = append(pds, pd)
			continue
		}

		// save message
		post.BeforeSave()

		if post.Type == model.MessageTypeIndicatorStart || post.Type == model.MessageTypeIndicatorEnd {
			messageIds = append(messageIds, post.MessageID)
			go publishMessage(ctx, post)
			continue
		}

		m, _ := datastore.Provider(ctx).SelectMessage(post.MessageID)
		if m != nil {
			pds = append(pds, &model.ProblemDetail{
				Message: "messageId is already exist",
				Status:  http.StatusConflict,
			})
		}

		lastMessage, err := datastore.Provider(ctx).InsertMessage(post)
		if err != nil {
			pd := &model.ProblemDetail{
				Message: "Message registration failed",
				Status:  http.StatusInternalServerError,
				Error:   errors.Wrap(err, ""),
			}
			pds = append(pds, pd)
			continue
		}

		logger.Debug(lastMessage)
		messageIds = append(messageIds, post.MessageID)

		// notification
		mi := &notification.MessageInfo{
			Text: fmt.Sprintf("[%s]%s", room.Name, lastMessage),
		}
		cfg := utils.Config()
		if cfg.Notification.DefaultBadgeCount != "" {
			dBadgeCount, err := strconv.Atoi(cfg.Notification.DefaultBadgeCount)
			if err == nil {
				mi.Badge = dBadgeCount
			}
		}
		go notification.Provider().Publish(ctx, room.NotificationTopicID, room.RoomID, mi)

		publishMessage(ctx, post)
		webhookMessage(ctx, post, user)
		logger.Info(fmt.Sprintf("Finish CreateMessage. Message=[%v]", post))
	}

	responseMessages := &model.ResponseMessages{
		MessageIds: messageIds,
		Errors:     pds,
	}
	for _, pd := range pds {
		logger.Error(pd.Error.Error())
	}
	return responseMessages
}

// GetMessage is get message
func GetMessage(ctx context.Context, messageID string) (*model.Message, *model.ProblemDetail) {
	if messageID == "" {
		return nil, &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*model.InvalidParam{
				&model.InvalidParam{
					Name:   "messageId",
					Reason: "messageId is required, but it's empty.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Create message failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if message == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}

	return message, nil
}

func publishMessage(ctx context.Context, message *model.Message) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(message.RoomID, datastore.WithRoleIDs([]int32{message.Role}))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(message)
	rtmEvent := &pbroker.RTMEvent{
		Type:    pbroker.MessageEvent,
		Payload: buffer.Bytes(),
		UserIDs: userIDs,
	}
	err = pbroker.Provider().PublishMessage(rtmEvent)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
