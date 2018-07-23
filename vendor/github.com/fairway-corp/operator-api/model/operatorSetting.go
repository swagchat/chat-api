package model

import (
	"time"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/utils"
)

type OperatorSetting struct {
	chatpb.OperatorSetting
}

func (os *OperatorSetting) ConvertToPbOperatorSetting() *chatpb.OperatorSetting {
	pbOs := &chatpb.OperatorSetting{}
	pbOs.SettingID = os.SettingID
	pbOs.SiteID = os.SiteID
	pbOs.Domain = os.Domain
	pbOs.SystemUserID = os.SystemUserID
	pbOs.FirstMessage = os.FirstMessage
	pbOs.TimeoutMessage = os.TimeoutMessage
	pbOs.OperatorBaseURL = os.OperatorBaseURL
	pbOs.NotificationSlackURL = os.NotificationSlackURL
	pbOs.Created = os.Created
	pbOs.Modified = os.Modified
	return pbOs
}

type CreateOperatorSettingRequest struct {
	chatpb.CreateOperatorSettingRequest
}

func (cgsr *CreateOperatorSettingRequest) Validate() *ErrorResponse {
	return nil
}

func (cosr *CreateOperatorSettingRequest) GenerateOperatorSetting() *OperatorSetting {
	nowTimestamp := time.Now().Unix()

	os := &OperatorSetting{}
	os.SettingID = utils.GenerateUUID()
	os.SiteID = cosr.SiteID
	os.Domain = cosr.Domain
	os.SystemUserID = cosr.SystemUserID
	os.FirstMessage = cosr.FirstMessage
	os.TimeoutMessage = cosr.TimeoutMessage
	os.OperatorBaseURL = cosr.OperatorBaseURL
	os.NotificationSlackURL = cosr.NotificationSlackURL
	os.Created = nowTimestamp
	os.Modified = nowTimestamp
	return os
}

type GetOperatorSettingRequest struct {
	chatpb.GetOperatorSettingRequest
}

func (ggsr *GetOperatorSettingRequest) Validate() *ErrorResponse {
	return nil
}

type UpdateOperatorSettingRequest struct {
	chatpb.UpdateOperatorSettingRequest
}

func (ugsr *UpdateOperatorSettingRequest) Validate() *ErrorResponse {
	return nil
}

func (uosr *UpdateOperatorSettingRequest) GenerateOperatorSetting() *OperatorSetting {
	nowTimestamp := time.Now().Unix()
	os := &OperatorSetting{}
	if uosr.SiteID != nil {
		os.SiteID = *uosr.SiteID
	}
	if uosr.Domain != nil {
		os.Domain = *uosr.Domain
	}
	if uosr.SystemUserID != nil {
		os.SystemUserID = *uosr.SystemUserID
	}
	if uosr.FirstMessage != nil {
		os.FirstMessage = uosr.FirstMessage
	}
	if uosr.FirstMessage != nil {
		os.TimeoutMessage = uosr.TimeoutMessage
	}
	if uosr.OperatorBaseURL != nil {
		os.OperatorBaseURL = *uosr.OperatorBaseURL
	}
	if uosr.NotificationSlackURL != nil {
		os.NotificationSlackURL = *uosr.NotificationSlackURL
	}
	os.Modified = nowTimestamp
	return os
}
