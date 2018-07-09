package model

const (
	ErrorInvalidParams         = "Invalid params"
	ErrorDatabase              = "Database error"
	ErrorNotification          = "Notification error"
	ErrorOperationNotPermitted = "Operation not permitted"
	ErrorUnauthorized          = "Unauthorized"
)

type ProblemDetail struct {
	Title         string         `json:"title,omitempty"`
	Message       string         `json:"message,omitempty"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
	Status        int            `json:"-"`
	Error         error          `json:"-"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// func (pd *ProblemDetail) Reset()         { *pd = ProblemDetail{} }
// func (pd *ProblemDetail) String() string { return "problem detail string" }
// func (pd *ProblemDetail) ProtoMessage()  {}
