// Business Logic

package services

import (
	"context"
	"net/http"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
)

func PostUser(post *models.User, jwt *models.JWT) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}
	post.BeforePost(jwt)

	user, err := datastore.Provider().SelectUser(post.UserId, true, true, true)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if user != nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource already exists",
			Status: http.StatusConflict,
		}
	}

	user, err = datastore.Provider().InsertUser(post)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	return user, nil
}

func GetUsers() (*models.Users, *models.ProblemDetail) {
	users, err := datastore.Provider().SelectUsers()
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.Users{
		Users: users,
	}, nil
}

func GetUser(userId string) (*models.User, *models.ProblemDetail) {
	user, err := datastore.Provider().SelectUser(userId, true, true, true)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

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

	user, err := datastore.Provider().UpdateUser(user)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Update user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return user, nil
}

func DeleteUser(userId string) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return pd
	}

	devices, err := datastore.Provider().SelectDevicesByUserId(userId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return pd
	}
	if devices != nil {
		for _, device := range devices {
			nRes := <-notification.Provider().DeleteEndpoint(device.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nRes.ProblemDetail
			}
		}
	}

	err = datastore.Provider().UpdateUserDeleted(userId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return pd
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByUserId(ctx, userId)

	return nil
}

func selectUser(userId string) (*models.User, *models.ProblemDetail) {
	user, err := datastore.Provider().SelectUser(userId, false, false, false)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return user, nil
}

func unsubscribeByUserId(ctx context.Context, userId string) {
	subscriptions, err := datastore.Provider().SelectDeletedSubscriptionsByUserId(userId)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	unsubscribe(ctx, subscriptions)
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
			Title:  "You do not have permission",
			Status: http.StatusUnauthorized,
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
			Title:  "You do not have permission",
			Status: http.StatusUnauthorized,
		}
	}

	return nil
}
