package model

import (
	scpb "github.com/swagchat/protobuf"
)

type errorResponseOptions struct {
	err           error
	invalidParams []*scpb.InvalidParam
}

type ErrorResponseOption func(*errorResponseOptions)

func WithError(err error) ErrorResponseOption {
	return func(ops *errorResponseOptions) {
		ops.err = err
	}
}

func WithInvalidParams(invalidParams []*scpb.InvalidParam) ErrorResponseOption {
	return func(ops *errorResponseOptions) {
		ops.invalidParams = invalidParams
	}
}

type ErrorResponse struct {
	scpb.ErrorResponse
	// Status is a HTTP status
	Status int `json:"-"`
	// Error is a error struct
	Error error `json:"-"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, status int, opts ...ErrorResponseOption) *ErrorResponse {
	opt := errorResponseOptions{}
	for _, o := range opts {
		o(&opt)
	}

	errRes := &ErrorResponse{}
	errRes.Message = message
	errRes.Status = status
	errRes.InvalidParams = opt.invalidParams
	errRes.Error = opt.err
	return errRes
}
