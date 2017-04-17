package services

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
)

func CreateRoom(post *models.Room) (*models.Room, *models.ProblemDetail) {
	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}
	post.BeforeSave()

	dp := datastore.GetProvider()
	dRes := <-dp.RoomInsert(post)
	return dRes.Data.(*models.Room), dRes.ProblemDetail
}

func GetRooms(values url.Values) (*models.Rooms, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.RoomSelectAll()
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	rooms := &models.Rooms{
		Rooms: dRes.Data.([]*models.Room),
	}
	return rooms, nil
}

func GetRoom(roomId string) (*models.Room, *models.ProblemDetail) {
	pd := IsExistRoomId(roomId)
	if pd != nil {
		return nil, pd
	}

	room, pd := getRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomSelectUsersForRoom(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room.Users = dRes.Data.([]*models.UserForRoom)
	return room, nil
}

func PutRoom(roomId string, put *models.Room) (*models.Room, *models.ProblemDetail) {
	pd := IsExistRoomId(roomId)
	if pd != nil {
		return nil, pd
	}

	room, pd := getRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	room.Put(put)
	if pd := room.IsValid(); pd != nil {
		return nil, pd
	}
	room.BeforeSave()

	dp := datastore.GetProvider()
	dRes := <-dp.RoomUpdate(room)
	return dRes.Data.(*models.Room), dRes.ProblemDetail
}

func DeleteRoom(roomId string) (*models.ResponseRoomUser, *models.ProblemDetail) {
	pd := IsExistRoomId(roomId)
	if pd != nil {
		return nil, pd
	}

	room, pd := getRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	np := notification.GetProvider()
	if np != nil {
		if room.NotificationTopicId != nil {
			nRes := <-np.DeleteTopic(*room.NotificationTopicId)
			if nRes.ProblemDetail != nil {
				return nil, nRes.ProblemDetail
			}
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomUsersSelect(&roomId, nil)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	ruRes := deleteRoomUsers(dRes.Data.([]*models.RoomUser))

	room.NotificationTopicId = nil
	room.Deleted = time.Now().UnixNano()
	dRes = <-dp.RoomUpdate(room)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	return ruRes, nil
}

func GetRoomMessages(roomId string, requestParams url.Values) (*models.Messages, *models.ProblemDetail) {
	pd := IsExistRoomId(roomId)
	if pd != nil {
		return nil, pd
	}

	var err error
	limit := 10
	offset := 0
	if limitArray, ok := requestParams["limit"]; ok {
		limit, err = strconv.Atoi(limitArray[0])
		if err != nil {
			return nil, &models.ProblemDetail{
				Title:     "Request parameter error. (Get room's message list)",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "limit",
						Reason: "limit is required, but it's empty.",
					},
				},
			}
		}
	}
	if offsetArray, ok := requestParams["offset"]; ok {
		offset, err = strconv.Atoi(offsetArray[0])
		if err != nil {
			return nil, &models.ProblemDetail{
				Title:     "Request parameter error. (Get room's message list)",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "offset",
						Reason: "offset is required, but it's empty.",
					},
				},
			}
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.MessageSelectAll(roomId, limit, offset)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	messages := &models.Messages{
		Messages: dRes.Data.([]*models.Message),
	}

	dRes = <-dp.MessageCount(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	messages.AllCount = dRes.Data.(*models.Messages).AllCount
	return messages, nil
}

func IsExistRoomId(roomId string) *models.ProblemDetail {
	if roomId == "" {
		return &models.ProblemDetail{
			Title:     "Request parameter error. (Get room's message list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}
	return nil
}

func getRoom(roomId string) (*models.Room, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.RoomSelect(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.Room), nil
}
