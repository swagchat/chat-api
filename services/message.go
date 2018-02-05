package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/swagchat/chat-api/bots"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/rtm"
	"github.com/swagchat/chat-api/utils"
)

func PostMessage(posts *models.Messages) *models.ResponseMessages {
	log.Println("============= PostMessage ==============")
	messageIds := make([]string, 0)
	errors := make([]*models.ProblemDetail, 0)
	var lastMessage string
	for _, post := range posts.Messages {
		room, pd := selectRoom(post.RoomId)
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

		user, pd := selectUser(post.UserId)
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
			log.Printf("%#v\n", pd)
			errors = append(errors, pd)
			continue
		}

		post.BeforeSave()
		log.Printf("%#v\n", post)
		dRes := datastore.GetProvider().InsertMessage(post)
		if dRes.ProblemDetail != nil {
			errors = append(errors, dRes.ProblemDetail)
			continue
		}
		if dRes.Data == nil {
			lastMessage = ""
		} else {
			lastMessage = dRes.Data.(string)
		}
		messageIds = append(messageIds, post.MessageId)
		log.Printf("%#v\n", post)

		mi := &notification.MessageInfo{
			Text: utils.AppendStrings("[", room.Name, "]", lastMessage),
		}
		cfg := utils.GetConfig()
		if cfg.Notification.DefaultBadgeCount != "" {
			dBadgeCount, err := strconv.Atoi(cfg.Notification.DefaultBadgeCount)
			if err == nil {
				mi.Badge = dBadgeCount
			}
		}
		ctx, _ := context.WithCancel(context.Background())
		go notification.GetProvider().Publish(ctx, room.NotificationTopicId, room.RoomId, mi)
		go publishMessage(post)
		go postMessageToBotService(*user.IsBot, post)
	}

	responseMessages := &models.ResponseMessages{
		MessageIds: messageIds,
		Errors:     errors,
	}
	return responseMessages
}

func GetMessage(messageId string) (*models.Message, *models.ProblemDetail) {
	if messageId == "" {
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

	dRes := datastore.GetProvider().SelectMessage(messageId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.Message), nil
}

func publishMessage(m *models.Message) {
	m.EventName = "message"
	bytes, err := json.Marshal(m)
	if err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	mi := &rtm.MessagingInfo{
		Message: string(bytes),
	}
	err = rtm.GetMessagingProvider().PublishMessage(mi)
	if err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
}

func postMessageToBotService(isBot bool, m *models.Message) {
	dRes := datastore.GetProvider().SelectUsersForRoom(m.RoomId)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	if dRes.Data != nil {
		userForRooms := dRes.Data.([]*models.UserForRoom)
		for _, u := range userForRooms {
			if !isBot && *u.IsBot && m.UserId != u.UserId {
				dRes := datastore.GetProvider().SelectBot(u.UserId)
				if dRes.ProblemDetail != nil {
					pdBytes, _ := json.Marshal(dRes.ProblemDetail)
					utils.AppLogger.Error("",
						zap.String("problemDetail", string(pdBytes)),
						zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
					)
				}
				b := dRes.Data.(*models.Bot)

				var cm models.CognitiveMap
				log.Printf("%#v\n", string(b.Cognitive))
				json.Unmarshal(b.Cognitive, &cm)
				log.Printf("%#v\n", cm)
				log.Printf("%#v\n", cm.Text)

				var p bots.Provider
				var cred utils.JSONText
				switch m.Type {
				case "text":
					p = bots.GetProvider(cm.Text.Name)
					cred = cm.Text.Credencial
				case "image":
					p = bots.GetProvider(cm.Image.Name)
					cred = cm.Image.Credencial
				default:
					p = bots.GetProvider(cm.Text.Name)
					cred = cm.Text.Credencial
				}
				res := p.Post(m, b, cred)
				if res.ProblemDetail != nil {
					pdBytes, _ := json.Marshal(dRes.ProblemDetail)
					utils.AppLogger.Error("",
						zap.String("problemDetail", string(pdBytes)),
						zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
					)
				}
				PostMessage(res.Messages)
			}
		}
	}
}
