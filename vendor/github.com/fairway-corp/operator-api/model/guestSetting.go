package model

import (
	"github.com/fairway-corp/chatpb"
)

type GuestSetting struct {
	chatpb.GuestSetting
}

func (gs *GuestSetting) ConvertToPbGestSetting() *chatpb.GuestSetting {
	pbGs := &chatpb.GuestSetting{}
	pbGs.EnableWebchat = gs.EnableWebchat
	return pbGs
}

type CreateGuestSettingRequest struct {
	chatpb.CreateGuestSettingRequest
}

func (cgsr *CreateGuestSettingRequest) Validate() *ErrorResponse {
	return nil
}

func (cgsr *CreateGuestSettingRequest) GenerateGuestSetting() *GuestSetting {
	gs := &GuestSetting{}
	gs.EnableWebchat = cgsr.EnableWebchat
	return gs
}

type GetGuestSettingRequest struct {
	chatpb.GetGuestSettingRequest
}

func (ggsr *GetGuestSettingRequest) Validate() *ErrorResponse {
	return nil
}

type UpdateGuestSettingRequest struct {
	chatpb.UpdateGuestSettingRequest
}

func (ugsr *UpdateGuestSettingRequest) Validate() *ErrorResponse {
	return nil
}

func (ugsr *UpdateGuestSettingRequest) GenerateGuestSetting() *GuestSetting {
	gs := &GuestSetting{}
	if ugsr.EnableWebchat != nil {
		gs.EnableWebchat = ugsr.EnableWebchat
	}
	return gs
}
