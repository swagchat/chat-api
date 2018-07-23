package service

import (
	"context"
	"fmt"
	"net/http"

	gimei "github.com/betchi/go-gimei"
	"github.com/fairway-corp/operator-api/idp"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
)

func CreateGuest(ctx context.Context, req *model.CreateGuestRequest) (*model.User, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  CreateGuest. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	userID, token, err := idp.Provider(ctx).Create()
	if err != nil {
		errRes = model.NewErrorResponse(http.StatusInternalServerError, err)
		return nil, errRes
	}

	gimei := gimei.NewName()
	pbReq := req.GenerateToPbCreateUserRequest()
	pbReq.Name = fmt.Sprintf("%s(%s)(仮)", gimei.Kanji(), gimei.Katakana())
	pbReq.UserID = userID
	pbReq.RoleIDs = []int32{utils.RoleGeneral, utils.RoleGuest}

	c := scpb.NewUserServiceClient(getChatAPIGRPCConn())
	pbUser, err := c.CreateUser(context.Background(), pbReq)
	if err != nil {
		errRes := model.NewErrorResponse(http.StatusInternalServerError, err)
		errRes.Message = "Failed creating guest user"
		return nil, errRes
	}

	metaData := utils.JSONText{}
	err = metaData.UnmarshalJSON(pbUser.MetaData)
	if err != nil {
		logger.Error(err.Error())
	}

	user := &model.User{*pbUser, metaData}
	user.AccessToken = token
	logger.Info(fmt.Sprintf("Finish GetGuest. Response[%#v]", user))
	return user, nil
}

func GetGuest(ctx context.Context, req *model.GetGuestRequest) (*model.User, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetGuest. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	pbReq := req.GenerateToPbGetUserRequest()
	c := scpb.NewUserServiceClient(getChatAPIGRPCConn())
	pbUser, err := c.GetUser(context.Background(), pbReq)
	if err != nil {
		errRes := model.NewErrorResponse(http.StatusInternalServerError, err)
		errRes.Message = "Failed getting guest user"
		return nil, errRes
	}

	if pbUser == nil {
		return nil, model.NewErrorResponse(http.StatusNotFound, nil)
	}

	metaData := utils.JSONText{}
	err = metaData.UnmarshalJSON(pbUser.MetaData)
	if err != nil {
		logger.Error(err.Error())
	}

	token, err := idp.Provider(ctx).GetToken()
	if err != nil {
		errRes = model.NewErrorResponse(http.StatusInternalServerError, err)
		return nil, errRes
	}
	pbUser.AccessToken = token

	user := &model.User{*pbUser, metaData}
	logger.Info(fmt.Sprintf("Finish GetGuest. Response[%#v]", user))
	return user, nil
}
