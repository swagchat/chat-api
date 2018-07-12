package model

const (
	ErrorInvalidParams         = "Invalid params"
	ErrorDatabase              = "Database error"
	ErrorNotification          = "Notification error"
	ErrorOperationNotPermitted = "Operation not permitted"
	ErrorUnauthorized          = "Unauthorized"
)

type ProblemDetail struct {
	// Message is a error message
	Message string `json:"message,omitempty"`
	// DeveloperMessage is a error message for developer
	DeveloperMessage string `json:"developerMessage,omitempty"`
	// Info is a detail of error contents (This can be URL)
	Info          string          `json:"info,omitempty"`
	InvalidParams []*InvalidParam `json:"invalidParams,omitempty"`
	// Status is a HTTP status
	Status int `json:"-"`
	// Error is a error struct
	Error error `json:"-"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// func (pd *ProblemDetail) Reset()         { *pd = ProblemDetail{} }
// func (pd *ProblemDetail) String() string { return "problem detail string" }
// func (pd *ProblemDetail) ProtoMessage()  {}
