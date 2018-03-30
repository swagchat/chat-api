// Business Logic

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

func PostUser(post *models.User, jwt *models.JWT) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}
	post.BeforePost(jwt)

	dRes := datastore.DatastoreProvider().SelectUser(post.UserId, true, true, true)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data != nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusConflict,
		}
	}

	dRes = datastore.DatastoreProvider().InsertUser(post)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	return dRes.Data.(*models.User), nil
}

func GetUsers() (*models.Users, *models.ProblemDetail) {
	dRes := datastore.DatastoreProvider().SelectUsers()
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	users := &models.Users{
		Users: dRes.Data.([]*models.User),
	}
	return users, nil
}

func GetUser(userId string) (*models.User, *models.ProblemDetail) {
	dRes := datastore.DatastoreProvider().SelectUser(userId, true, true, true)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}

	user := dRes.Data.(*models.User)
	// unreadCountRooms := make([]*models.RoomForUser, 0)
	// notUnreadCountRooms := make([]*models.RoomForUser, 0)
	// for _, roomForUser := range user.Rooms {
	// 	if roomForUser.RuUnreadCount > 0 {
	// 		unreadCountRooms = append(unreadCountRooms, roomForUser)
	// 	} else {
	// 		notUnreadCountRooms = append(notUnreadCountRooms, roomForUser)
	// 	}
	// }
	// mergeRooms := append(unreadCountRooms, notUnreadCountRooms...)
	// user.Rooms = mergeRooms
	return user, nil
}

func GetProfile(userId string) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	return user, nil
}

func PutUser(put *models.User) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(put.UserId)
	if pd != nil {
		return nil, pd
	}

	if pd := user.IsValidPut(); pd != nil {
		return nil, pd
	}

	user.BeforePut(put)
	dRes := datastore.DatastoreProvider().UpdateUser(user)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	return user, nil
}

func DeleteUser(userId string) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return pd
	}

	dRes := datastore.DatastoreProvider().SelectDevicesByUserId(userId)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}
	if dRes.Data != nil {
		devices := dRes.Data.([]*models.Device)
		for _, device := range devices {
			nRes := <-notification.NotificationProvider().DeleteEndpoint(device.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nRes.ProblemDetail
			}
		}
	}

	dRes = datastore.DatastoreProvider().UpdateUserDeleted(userId)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByUserId(ctx, userId)

	return nil
}

func selectUser(userId string) (*models.User, *models.ProblemDetail) {
	dRes := datastore.DatastoreProvider().SelectUser(userId, false, false, false)
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

func unsubscribeByUserId(ctx context.Context, userId string) {
	dRes := datastore.DatastoreProvider().SelectDeletedSubscriptionsByUserId(userId)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	unsubscribe(ctx, dRes.Data.([]*models.Subscription))
}

func GetUserUnreadCount(userId string) (*models.UserUnreadCount, *models.ProblemDetail) {
	user, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	userUnreadCount := &models.UserUnreadCount{
		UnreadCount: user.UnreadCount,
	}
	return userUnreadCount, nil
}

func UserAuth(userId, sub string) *models.ProblemDetail {
	if userId != sub {
		return &models.ProblemDetail{
			Title:     "You do not have permission",
			Status:    http.StatusUnauthorized,
			ErrorName: models.ERROR_NAME_UNAUTHORIZED,
		}
	}

	return nil
}

func ContactsAuth(userId, sub string) *models.ProblemDetail {
	contacts, pd := GetContacts(sub)
	if pd != nil {
		return pd
	}

	isAuthorized := false
	for _, contact := range contacts.Users {
		if contact.UserId == userId {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return &models.ProblemDetail{
			Title:     "You do not have permission",
			Status:    http.StatusUnauthorized,
			ErrorName: models.ERROR_NAME_UNAUTHORIZED,
		}
	}

	return nil
}
