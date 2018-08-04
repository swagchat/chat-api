package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func confirmAssetExist(ctx context.Context, assetID string) (*model.Asset, *model.ErrorResponse) {
	asset, err := datastore.Provider(ctx).SelectAsset(assetID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if asset == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That asset is not exist. assetId[%s]", assetID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}
	return asset, nil
}

func confirmDeviceExist(ctx context.Context, userID string, platform scpb.Platform) (*model.Device, *model.ErrorResponse) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if device == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That device is not exist. userId[%s] platform[%d]", userID, platform),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return device, nil
}

func confirmDeviceNotExist(ctx context.Context, userID string, platform scpb.Platform) (*model.Device, *model.ErrorResponse) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if device != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That device already exist. userId[%s] platform[%d]", userID, platform),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return device, nil
}

func confirmMessageExist(ctx context.Context, messageID string) (*model.Message, *model.ErrorResponse) {
	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if message == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: fmt.Sprintf("That message is not exist. messageId[%s]", messageID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return message, nil
}

func confirmMessageNotExist(ctx context.Context, messageID string) (*model.Message, *model.ErrorResponse) {
	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if message != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: fmt.Sprintf("That message already exist. messageId[%s]", messageID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return message, nil
}

func confirmRoomExist(ctx context.Context, roomID string, opts ...datastore.SelectRoomOption) (*model.Room, *model.ErrorResponse) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if room == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: fmt.Sprintf("That room is not exist. roomId[%s]", roomID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return room, nil
}

func confirmRoomUserExist(ctx context.Context, roomID, userID string) (*model.RoomUser, *model.ErrorResponse) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUser(roomID, userID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if roomUser == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId | userId",
				Reason: fmt.Sprintf("That room user is not exist. roomId[%s] userId[%s]", roomID, userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return roomUser, nil
}

func confirmUserExist(ctx context.Context, userID string, opts ...datastore.SelectUserOption) (*model.User, *model.ErrorResponse) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if user == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: fmt.Sprintf("That user is not exist. userId[%s]", userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return user, nil
}

func confirmUserNotExist(ctx context.Context, userID string, opts ...datastore.SelectUserOption) (*model.User, *model.ErrorResponse) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if user != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: fmt.Sprintf("That user already exist. userId[%s]", userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return user, nil
}

func confirmUserIDsExist(ctx context.Context, requestUserIDs []string, keyName string) *model.ErrorResponse {
	existUserIDs, err := datastore.Provider(ctx).SelectUserIDsOfUser(requestUserIDs)
	if err != nil {
		return model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
	}

	if len(existUserIDs) != len(requestUserIDs) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   keyName,
				Reason: "It contains a userId that does not exist.",
			},
		}
		return model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return nil
}
