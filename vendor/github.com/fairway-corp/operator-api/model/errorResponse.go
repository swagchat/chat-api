package model

import (
	scpb "github.com/swagchat/protobuf"
)

type ErrorResponse struct {
	scpb.ErrorResponse
	// Status is a HTTP status
	Status int `json:"-"`
	// Error is a error struct
	Error error `json:"-"`
}

func NewErrorResponse(status int, err error) *ErrorResponse {
	errRes := &ErrorResponse{}
	errRes.Status = status
	errRes.Error = err
	return errRes
}
