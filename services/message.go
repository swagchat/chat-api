package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/bots"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/rtm"
	"github.com/swagchat/chat-api/utils"
)

// PostMessage is post message
func PostMessage(ctx context.Context, posts *models.Messages) *models.ResponseMessages {
	messageIds := make([]string, 0)
	errors := make([]*models.ProblemDetail, 0)
	for _, post := range posts.Messages {
		room, pd := selectRoom(ctx, post.RoomID)
		if pd != nil {
			errors = append(errors, &models.ProblemDetail{
				Title:     "Request parameter error. (Create message item)",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "roomId",
						Reason: "roomId is invalid. Not exist room.",
					},
				},
			})
			continue
		}

		user, pd := selectUser(ctx, post.UserID, datastore.WithRoles(true))
		if pd != nil {
			errors = append(errors, &models.ProblemDetail{
				Title:     "Request parameter error. (Create message item)",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "userId",
						Reason: "userId is invalid. Not exist user.",
					},
				},
			})
			continue
		}

		if pd := post.IsValid(); pd != nil {
			errors = append(errors, pd)
			continue
		}

		if post.Type == models.MessageTypeIndicatorStart || post.Type == models.MessageTypeIndicatorEnd {
			messageID := ""
			switch post.Type {
			case models.MessageTypeIndicatorStart:
				messageID = fmt.Sprintf("%s-%s", models.MessageTypeIndicatorStart, post.UserID)
			case models.MessageTypeIndicatorEnd:
				messageID = fmt.Sprintf("%s-%s", models.MessageTypeIndicatorEnd, post.UserID)
			}
			messageIds = append(messageIds, messageID)
			post.MessageID = messageID
			post.Created = time.Now().Unix()
			go rtmPublish(ctx, post)
			continue
		}

		// save message
		post.BeforeSave()
		lastMessage, err := datastore.Provider(ctx).InsertMessage(post)
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Message registration failed",
				Status: http.StatusInternalServerError,
				Error:  err,
			}
			errors = append(errors, pd)
			continue
		}
		messageIds = append(messageIds, post.MessageID)

		// notification
		mi := &notification.MessageInfo{
			Text: utils.AppendStrings("[", room.Name, "]", lastMessage),
		}
		cfg := utils.Config()
		if cfg.Notification.DefaultBadgeCount != "" {
			dBadgeCount, err := strconv.Atoi(cfg.Notification.DefaultBadgeCount)
			if err == nil {
				mi.Badge = dBadgeCount
			}
		}
		go notification.Provider().Publish(ctx, room.NotificationTopicID, room.RoomID, mi)

		rtmPublish(ctx, post)
		postMessageToBotService(ctx, post, user)
		webhookMessage(ctx, post, user)
	}

	responseMessages := &models.ResponseMessages{
		MessageIds: messageIds,
		Errors:     errors,
	}
	return responseMessages
}

// GetMessage is get message
func GetMessage(ctx context.Context, messageID string) (*models.Message, *models.ProblemDetail) {
	if messageID == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Get message item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "messageId",
					Reason: "messageId is required, but it's empty.",
				},
			},
		}
	}

	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if message == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return message, nil
}

func postMessageToBotService(ctx context.Context, m *models.Message, u *models.User) {
	if u.IsRole(models.RoleOperator) || u.IsRole(models.RoleBot) {
		return
	}

	usersForRoom, err := datastore.Provider(ctx).SelectUsersForRoom(m.RoomID)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	if len(usersForRoom) > 0 {
		for _, ufr := range usersForRoom {
			if *ufr.IsBot && m.UserID != ufr.UserID {
				bot, err := datastore.Provider(ctx).SelectBot(ufr.UserID)
				if err != nil {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Error: err,
					})
				}

				var cm models.CognitiveMap
				json.Unmarshal(bot.Cognitive, &cm)

				var res *bots.BotResult
				switch m.Type {
				case "text":
					p := bots.Provider(cm.Text.Name)
					cred := cm.Text.Credencial
					res = p.Post(m, bot, cred)
				case "image":
					continue
				default:
					continue
				}
				if res != nil {
					PostMessage(ctx, res.Messages)
				}
			}
		}
	}
}

func rtmPublish(ctx context.Context, message *models.Message) {
	var roles []models.Role
	if message.SuggestMessageID != "" {
		roles = []models.Role{models.RoleOperator}
	}

	userIDs, err := datastore.Provider(ctx).SelectRoomUserIDsByRoomID(message.RoomID, roles)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(message)
	rtmEvent := &rtm.RTMEvent{
		Type:    rtm.MessageEvent,
		Payload: buffer.Bytes(),
		UserIDs: userIDs,
	}
	err = rtm.Provider().Publish(rtmEvent)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
}
