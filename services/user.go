// Business Logic

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateUser(post *models.User) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}
	post.BeforeSave()

	dRes := <-datastore.GetProvider().UserInsert(post)
	return dRes.Data.(*models.User), dRes.ProblemDetail
}

func GetUsers() (*models.Users, *models.ProblemDetail) {
	dRes := <-datastore.GetProvider().UserSelectAll()
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	users := &models.Users{
		Users: dRes.Data.([]*models.User),
	}
	return users, nil
}

func GetUser(userId string) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	dRes := <-datastore.GetProvider().UserSelectRoomsForUser(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	user.Rooms = dRes.Data.([]*models.RoomForUser)
	return user, pd
}

func PutUser(userId string, put *models.User) (*models.User, *models.ProblemDetail) {
	user, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	user.Put(put)
	if pd := user.IsValid(); pd != nil {
		return nil, pd
	}
	user.BeforeSave()

	if *user.UnreadCount == 0 {
		dRes := <-datastore.GetProvider().RoomUserMarkAllAsRead(userId)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
	}

	dRes := <-datastore.GetProvider().UserUpdate(user)
	return dRes.Data.(*models.User), dRes.ProblemDetail
}

func DeleteUser(userId string) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return pd
	}

	dRes := <-datastore.GetProvider().DeviceSelectByUserId(userId)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}
	if dRes.Data != nil {
		devices := dRes.Data.([]*models.Device)
		for _, device := range devices {
			nRes := <-notification.GetProvider().DeleteEndpoint(device.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nRes.ProblemDetail
			}
		}
	}

	dRes = <-datastore.GetProvider().UserUpdateDeleted(userId)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go unsubscribeByUserId(ctx, userId)

	return nil
}

func selectUser(userId string) (*models.User, *models.ProblemDetail) {
	dRes := <-datastore.GetProvider().UserSelect(userId, false, false)
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
	dRes := <-datastore.GetProvider().SubscriptionSelectByUserId(userId)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	unsubscribe(ctx, dRes.Data.([]*models.Subscription))
}
