package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

// ContactsAuthz is contacts authorize
func ContactsAuthz(ctx context.Context, requestUserID, resourceUserID string) *model.ErrorResponse {
	req := &model.GetContactsRequest{}
	req.UserID = requestUserID

	contacts, errRes := GetContacts(ctx, req)
	if errRes != nil {
		return errRes
	}

	isAuthorized := false
	for _, contact := range contacts.Users {
		if contact.UserID == resourceUserID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return model.NewErrorResponse("You do not have permission", http.StatusUnauthorized)
	}

	return nil
}

// RoomAuthz is room authorize
func RoomAuthz(ctx context.Context, roomID, userID string) *model.ErrorResponse {
	room, errRes := confirmRoomExist(ctx, roomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		return errRes
	}

	if room.Type == scpb.RoomType_RoomTypePublicRoom {
		return nil
	}

	isAuthorized := false
	for _, user := range room.Users {
		if user.UserID == userID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return model.NewErrorResponse("You are not this room member", http.StatusUnauthorized)
	}

	return nil
}
