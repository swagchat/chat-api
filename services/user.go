// Business Logic

package services

import (
	"log"
	"net/http"
	"time"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateUser(requestUser *models.User) (*models.User, *models.ProblemDetail) {
	if requestUser.Name == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create user item)",
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

	userId := requestUser.UserId
	if userId != "" && !utils.IsValidId(userId) {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}
	if userId == "" {
		userId = utils.CreateUuid()
	}

	var metaData []byte
	if requestUser.MetaData == nil {
		metaData = []byte("{}")
	} else {
		metaData = requestUser.MetaData
	}

	unreadCount := int64(0)

	user := &models.User{
		UserId:         userId,
		Name:           requestUser.Name,
		PictureUrl:     requestUser.PictureUrl,
		InformationUrl: requestUser.InformationUrl,
		UnreadCount:    &unreadCount,
		MetaData:       metaData,
		Created:        time.Now().UnixNano(),
		Modified:       time.Now().UnixNano(),
	}
	if requestUser.DeviceToken != nil {
		if *requestUser.DeviceToken == "" {
			user.DeviceToken = nil
		} else {
			user.DeviceToken = requestUser.DeviceToken
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.UserInsert(user)
	return dRes.Data.(*models.User), dRes.ProblemDetail
}

func GetUsers() (*models.Users, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.UserSelectAll()
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	users := &models.Users{
		Users: dRes.Data.([]*models.User),
	}
	return users, nil
}

func GetUser(userId string) (*models.User, *models.ProblemDetail) {
	if userId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Update user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.UserSelect(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	user := dRes.Data.(*models.User)
	dRes = <-dp.UserSelectRoomsForUser(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	user.Rooms = dRes.Data.([]*models.RoomForUser)
	return user, nil
}

func PutUser(userId string, requestUser *models.User) (*models.User, *models.ProblemDetail) {
	if userId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Update user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.UserSelect(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	user := dRes.Data.(*models.User)

	if requestUser.DeviceToken != nil {
		np := notification.GetProvider()
		log.Println("np", np)
		if np != nil {
			if user.DeviceToken != nil && *requestUser.DeviceToken == "" {
				if user.NotificationDeviceId != nil {
					nRes := <-np.DeleteEndpoint(*user.NotificationDeviceId)
					if nRes.ProblemDetail != nil {
						return nil, nRes.ProblemDetail
					}
					user.DeviceToken = nil
					user.NotificationDeviceId = nil
					user.Modified = time.Now().UnixNano()
					dRes := <-dp.UserUpdate(user)
					if dRes.ProblemDetail != nil {
						return nil, dRes.ProblemDetail
					}
				}

				dRes := <-dp.RoomUsersSelect(nil, []string{userId})
				if dRes.ProblemDetail != nil {
					return nil, dRes.ProblemDetail
				}
				ruRes := unsubscribeAllRoom(dRes.Data.([]*models.RoomUser))
				if len(ruRes.Errors) > 0 {
					return nil, &models.ProblemDetail{
						Title:     "Updating user item error. (Update user item)",
						Status:    http.StatusInternalServerError,
						ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
					}
				}
				user.DeviceToken = nil
			} else if (user.DeviceToken == nil && *requestUser.DeviceToken != "") ||
				(user.DeviceToken != nil && (user.DeviceToken != requestUser.DeviceToken)) {

				nRes := <-np.CreateEndpoint(*requestUser.DeviceToken)
				if nRes.ProblemDetail != nil {
					return nil, &models.ProblemDetail{
						Title:     "Updating user item error. (Update user item)",
						Status:    http.StatusInternalServerError,
						ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
					}
				}
				if nRes.Data == nil {
					return nil, &models.ProblemDetail{
						Title:     "Creating notification endpoint. (Update user item)",
						Status:    http.StatusInternalServerError,
						ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
					}
				}
				notificationDeviceId := nRes.Data.(*string)

				dRes := <-dp.RoomUsersSelect(nil, []string{userId})
				if dRes.ProblemDetail != nil {
					return nil, dRes.ProblemDetail
				}

				ruRes := subscribeAllRoom(dRes.Data.([]*models.RoomUser), *notificationDeviceId)
				if len(ruRes.Errors) > 0 {
					return nil, &models.ProblemDetail{
						Title:     "Updating user item error. (Update user item)",
						Status:    http.StatusInternalServerError,
						ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
					}
				}
				user.DeviceToken = requestUser.DeviceToken
				user.NotificationDeviceId = notificationDeviceId
			}
		}
	}

	if requestUser.Name != "" {
		user.Name = requestUser.Name
	}
	if requestUser.PictureUrl != "" {
		user.PictureUrl = requestUser.PictureUrl
	}
	if requestUser.InformationUrl != "" {
		user.InformationUrl = requestUser.InformationUrl
	}
	if requestUser.UnreadCount != nil {
		user.UnreadCount = requestUser.UnreadCount
		var zero int64
		zero = 0
		if requestUser.UnreadCount == &zero {
			dRes := <-dp.RoomUserMarkAllAsRead(userId)
			if dRes.ProblemDetail != nil {
				return nil, dRes.ProblemDetail
			}
		}
	}
	if requestUser.MetaData != nil {
		user.MetaData = requestUser.MetaData
	}
	user.Modified = time.Now().UnixNano()

	dRes = <-dp.UserUpdate(user)
	return dRes.Data.(*models.User), dRes.ProblemDetail
}

func DeleteUser(userId string) (*models.ResponseRoomUser, *models.ProblemDetail) {
	if userId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Delete user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.UserSelect(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	user := dRes.Data.(*models.User)

	np := notification.GetProvider()
	if np != nil {
		if user.NotificationDeviceId != nil {
			nRes := <-np.DeleteEndpoint(*user.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nil, nRes.ProblemDetail
			}
		}
	}

	dRes = <-dp.RoomUsersSelect(nil, []string{userId})
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	ruRes := deleteRoomUsers(dRes.Data.([]*models.RoomUser))

	user.DeviceToken = nil
	user.NotificationDeviceId = nil
	user.Deleted = time.Now().UnixNano()
	dRes = <-dp.UserUpdate(user)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	return ruRes, nil
}

/*
func GetUserRooms(userId string) (*models.User, *models.ProblemDetail) {
	if userId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter for user's room list getting is invalid.",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.UserSelectUserRooms(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	user := &models.User{
		Rooms: dRes.Data.([]*models.UserRoom),
	}

	return user, nil
}
*/
