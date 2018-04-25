package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

		user, pd := selectUser(ctx, post.UserID)
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
		go postMessageToBotService(ctx, *user.IsBot, post)
		go func() {
			userIDs, err := datastore.Provider(ctx).SelectRoomUserIDsByRoomID(room.RoomID)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: err,
				})
			}

			buffer := new(bytes.Buffer)
			json.NewEncoder(buffer).Encode(post)
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
		}()
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

func postMessageToBotService(ctx context.Context, isBot bool, m *models.Message) {
	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(m.RoomID)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	if len(userForRooms) > 0 {
		for _, u := range userForRooms {
			if !isBot && *u.IsBot && m.UserID != u.UserID {
				bot, err := datastore.Provider(ctx).SelectBot(u.UserID)
				if err != nil {
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						Error: err,
					})
				}

				var cm models.CognitiveMap
				json.Unmarshal(bot.Cognitive, &cm)

				var res bots.BotResult
				switch m.Type {
				case "text":
					p := bots.Provider(cm.Text.Name)
					cred := cm.Text.Credencial
					res = p.Post(m, bot, cred)
				case "image":
					p := bots.Provider(cm.Image.Name)
					cred := cm.Image.Credencial
					res = p.Post(m, bot, cred)
				default:
					p := bots.Provider(cm.Text.Name)
					cred := cm.Text.Credencial
					res = p.Post(m, bot, cred)
				}
				PostMessage(ctx, res.Messages)
			}
		}
	}
}
