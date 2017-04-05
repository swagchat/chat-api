package services

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/messaging"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateMessage(requestMessages *models.Messages) (*models.ResponseMessages, *models.ProblemDetail) {
	// TODO transaction
	messageIds := make([]string, len(requestMessages.Messages))
	var messageId string
	var lastMessage string
	for i, requestMessage := range requestMessages.Messages {
		messageId = requestMessage.MessageId
		if messageId != "" && !utils.IsValidId(messageId) {
			return nil, &models.ProblemDetail{
				Title:     "Request parameter error. (Create message item)",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "messageId",
						Reason: "messageId is invalid. Available characters are alphabets, numbers and hyphens.",
					},
				},
			}
		}
		if messageId == "" {
			messageId = utils.CreateUuid()
		}

		messageIds[i] = messageId
		message := &models.Message{
			RoomId:    requestMessage.RoomId,
			UserId:    requestMessage.UserId,
			MessageId: messageId,
			Type:      requestMessage.Type,
			Payload:   requestMessage.Payload,
			Created:   time.Now().UnixNano(),
			Modified:  time.Now().UnixNano(),
		}
		dp := datastore.GetProvider()
		dRes := <-dp.MessageInsert(message)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}

		switch requestMessage.Type {
		case "text":
			var payloadText models.PayloadText
			json.Unmarshal(message.Payload, &payloadText)
			log.Printf("%#v\n", payloadText)
			lastMessage = payloadText.Text
		case "image":
			lastMessage = "画像を送信しました"
		case "location":
			lastMessage = "位置情報を送信しました"
		}

		dRes = <-dp.RoomSelect(requestMessage.RoomId)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}

		room := dRes.Data.(*models.Room)
		room.LastMessage = lastMessage
		room.LastMessageUpdated = time.Now().UnixNano()
		dRes = <-dp.RoomUpdate(room)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}

		dRes = <-dp.RoomUserUnreadCountUp(requestMessage.RoomId, requestMessage.UserId)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}

		dRes = <-dp.RoomUserUsersSelect(requestMessage.RoomId)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
		for _, user := range dRes.Data.([]*models.User) {
			if user.UserId == requestMessage.UserId {
				continue
			}
			dRes = <-dp.UserUnreadCountUp(user.UserId)
			if dRes.ProblemDetail != nil {
				return nil, dRes.ProblemDetail
			}
		}

		if utils.Cfg.RealtimeServer.Endpoint != "" {
			go func() {
				message.EventName = "message"
				bytes, err := json.Marshal(message)
				if err != nil {
					utils.AppLogger.Error("",
						zap.String("msg", err.Error()),
					)
				}
				messagingInfo := &messaging.MessagingInfo{
					Message: string(bytes),
				}
				messagingProvider := messaging.GetMessagingProvider()
				err = messagingProvider.PublishMessage(messagingInfo)
				if err != nil {
					utils.AppLogger.Error("",
						zap.String("msg", err.Error()),
					)
				}
			}()
		}

		np := notification.GetProvider()
		if np != nil {
			dp := datastore.GetProvider()
			log.Println(requestMessage.RoomId)
			log.Println(requestMessage.UserId)
			dRes := <-dp.RoomUserSelect(requestMessage.RoomId, requestMessage.UserId)
			roomUser := dRes.Data.(*models.RoomUser)
			log.Println(roomUser)
			messageInfo := &notification.MessageInfo{
				Text:  utils.AppendStrings("[", room.Name, "]", lastMessage),
				Badge: *roomUser.UnreadCount,
			}
			nRes := <-np.Publish(*room.NotificationTopicId, messageInfo)
			if nRes.ProblemDetail != nil {
				problemDetailBytes, _ := json.Marshal(nRes.ProblemDetail)
				utils.AppLogger.Error("",
					zap.String("problemDetail", string(problemDetailBytes)),
				)
			}
		}
	}

	responseMessages := &models.ResponseMessages{
		MessageIds: messageIds,
	}
	return responseMessages, nil
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

	dp := datastore.GetProvider()
	dRes := <-dp.MessageSelect(messageId)
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
