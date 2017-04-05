package services

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateRoom(requestRoom *models.Room) (*models.Room, *models.ProblemDetail) {
	if requestRoom.Name == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "name",
					Reason: "name is required, but it's empty.",
				},
			},
		}
	}

	roomId := requestRoom.RoomId
	if roomId != "" && !utils.IsValidId(roomId) {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}
	if roomId == "" {
		roomId = utils.CreateUuid()
	}

	var metaData []byte
	if requestRoom.MetaData == nil {
		metaData = []byte("{}")
	} else {
		metaData = requestRoom.MetaData
	}

	var isPublic bool
	if requestRoom.IsPublic == nil {
		isPublic = false
	} else {
		isPublic = *requestRoom.IsPublic
	}

	room := &models.Room{
		RoomId:         roomId,
		Name:           requestRoom.Name,
		PictureUrl:     requestRoom.PictureUrl,
		InformationUrl: requestRoom.InformationUrl,
		MetaData:       metaData,
		IsPublic:       &isPublic,
		Created:        time.Now().UnixNano(),
		Modified:       time.Now().UnixNano(),
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomInsert(room)
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
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Get room item)",
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
	room := dRes.Data.(*models.Room)
	dRes = <-dp.RoomSelectUsersForRoom(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room.Users = dRes.Data.([]*models.UserForRoom)
	return room, nil
}

func PutRoom(roomId string, requestRoom *models.Room) (*models.Room, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Update room item)",
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
	room := dRes.Data.(*models.Room)

	if requestRoom.Name != "" {
		room.Name = requestRoom.Name
	}
	if requestRoom.PictureUrl != "" {
		room.PictureUrl = requestRoom.PictureUrl
	}
	if requestRoom.InformationUrl != "" {
		room.InformationUrl = requestRoom.InformationUrl
	}
	if requestRoom.MetaData != nil {
		room.MetaData = requestRoom.MetaData
	}
	if requestRoom.IsPublic != nil {
		room.IsPublic = requestRoom.IsPublic
	}
	room.Modified = time.Now().UnixNano()

	dRes = <-dp.RoomUpdate(room)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room = dRes.Data.(*models.Room)

	return room, nil
}

func DeleteRoom(roomId string) (*models.ResponseRoomUser, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Delete room item)",
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
	room := dRes.Data.(*models.Room)

	np := notification.GetProvider()
	if np != nil {
		if room.NotificationTopicId != nil {
			nRes := <-np.DeleteTopic(*room.NotificationTopicId)
			if nRes.ProblemDetail != nil {
				return nil, nRes.ProblemDetail
			}
		}
	}

	dRes = <-dp.RoomUsersSelect(&roomId, nil)
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
	if roomId == "" {
		return nil, &models.ProblemDetail{
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
