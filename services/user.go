// Business Logic

package services

import (
	"net/http"
	"time"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
)

func CreateUser(post *models.User) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}
	post.BeforeSave()

	dp := datastore.GetProvider()
	dRes := <-dp.UserInsert(post)
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
	pd := IsExistUserId(userId)
	if pd != nil {
		return nil, pd
	}

	user, pd := getUser(userId)
	return user, pd
}

func PutUser(userId string, put *models.User) (*models.User, *models.ProblemDetail) {
	pd := IsExistUserId(userId)
	if pd != nil {
		return nil, pd
	}

	user, pd := getUser(userId)
	if pd != nil {
		return nil, pd
	}

	user.Put(put)
	if pd := user.IsValid(); pd != nil {
		return nil, pd
	}
	user.BeforeSave()

	dp := datastore.GetProvider()
	if *user.UnreadCount == 0 {
		dRes := <-dp.RoomUserMarkAllAsRead(userId)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
	}

	dRes := <-dp.UserUpdate(user)
	return dRes.Data.(*models.User), dRes.ProblemDetail
}

func DeleteUser(userId string) *models.ProblemDetail {
	pd := IsExistUserId(userId)
	if pd != nil {
		return pd
	}

	user, pd := getUser(userId)
	if pd != nil {
		return pd
	}
	//	np := notification.GetProvider()
	//	if user.Devices != nil {
	//		for _, device := range user.Devices {
	//			if device.NotificationDeviceId != nil {
	//				nRes := <-np.DeleteEndpoint(*device.NotificationDeviceId)
	//				if nRes.ProblemDetail != nil {
	//					return nRes.ProblemDetail
	//				}
	//			}
	//			dRes := <-dp.DeviceDelete(user.UserId, device.Platform)
	//			if dRes.ProblemDetail != nil {
	//				return dRes.ProblemDetail
	//			}
	//		}
	//	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomUsersSelect(nil, []string{userId})
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}
	go deleteRoomUsers(dRes.Data.([]*models.RoomUser))

	user.Deleted = time.Now().UnixNano()
	dRes = <-dp.UserUpdate(user)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	return nil
}

func IsExistUserId(userId string) *models.ProblemDetail {
	if userId == "" {
		return &models.ProblemDetail{
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
	return nil
}

func getUser(userId string) (*models.User, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.UserSelect(userId, false, false)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.User), nil
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
