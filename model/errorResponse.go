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

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, invalidParams []*scpb.InvalidParam, status int, err error) *ErrorResponse {
	errRes := &ErrorResponse{}
	errRes.Message = message
	errRes.InvalidParams = invalidParams
	errRes.Status = status
	errRes.Error = err
	return errRes
}
