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

// PostUser is post user
func PostUser(ctx context.Context, post *models.User) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}

	post.BeforePost()

	user, err := datastore.Provider(ctx).SelectUser(post.UserID, datastore.WithBlocks(true), datastore.WithDevices(true), datastore.WithRooms(true))
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

	user, err = datastore.Provider(ctx).InsertUser(post)
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

// GetUsers is get users
func GetUsers(ctx context.Context) (*models.Users, *models.ProblemDetail) {
	users, err := datastore.Provider(ctx).SelectUsers()
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

// GetUser is get user
func GetUser(ctx context.Context, userID string) (*models.User, *models.ProblemDetail) {
	user, err := datastore.Provider(ctx).SelectUser(userID, datastore.WithBlocks(true), datastore.WithDevices(true), datastore.WithRooms(true))
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

// PutUser is put user
func PutUser(ctx context.Context, put *models.User) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(ctx, put.UserID)
	if pd != nil {
		return nil, pd
	}

	if pd := user.IsValidPut(); pd != nil {
		return nil, pd
	}

	user.BeforePut(put)

	user, err := datastore.Provider(ctx).UpdateUser(user)
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

// DeleteUser is delete user
func DeleteUser(ctx context.Context, userID string) *models.ProblemDetail {
	dsp := datastore.Provider(ctx)
	// User existence check
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return pd
	}

	devices, err := dsp.SelectDevicesByUserID(userID)
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
			nRes := <-notification.Provider().DeleteEndpoint(device.NotificationDeviceID)
			if nRes.ProblemDetail != nil {
				return nRes.ProblemDetail
			}
		}
	}

	err = dsp.UpdateUserDeleted(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return pd
	}

	go unsubscribeByUserID(ctx, userID)

	return nil
}

// GetUserUnreadCount is get user unread count
func GetUserUnreadCount(ctx context.Context, userID string) (*models.UserUnreadCount, *models.ProblemDetail) {
	user, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	userUnreadCount := &models.UserUnreadCount{
		UnreadCount: user.UnreadCount,
	}
	return userUnreadCount, nil
}

// GetContacts is get contacts
func GetContacts(ctx context.Context, userID string) (*models.Users, *models.ProblemDetail) {
	contacts, err := datastore.Provider(ctx).SelectContacts(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get contact list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.Users{
		Users: contacts,
	}, nil
}

// GetProfile is get profile
func GetProfile(ctx context.Context, userID string) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	return user, nil
}

func selectUser(ctx context.Context, userID string, opts ...interface{}) (*models.User, *models.ProblemDetail) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
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

func unsubscribeByUserID(ctx context.Context, userID string) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByUserID(userID)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	unsubscribe(ctx, subscriptions)
}

// ContactsAuthz is contacts authorize
func ContactsAuthz(ctx context.Context, requestUserID, resourceUserID string) *models.ProblemDetail {
	contacts, pd := GetContacts(ctx, requestUserID)
	if pd != nil {
		return pd
	}

	isAuthorized := false
	for _, contact := range contacts.Users {
		if contact.UserID == resourceUserID {
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
