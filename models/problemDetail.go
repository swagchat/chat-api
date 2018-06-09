package models

const (
	ERROR_NAME_INVALID_JSON            = "invalid-json"
	ERROR_NAME_INVALID_PARAM           = "invalid-param"
	ERROR_NAME_DATABASE_ERROR          = "database-error"
	ERROR_NAME_NOTIFICATION_ERROR      = "notification-error"
	ERROR_NAME_OPERATION_NOT_PERMITTED = "operation-not-permitted"
	ERROR_NAME_UNAUTHORIZED            = "unauthorized"
)

type ProblemDetail struct {
	Type          string         `json:"type,omitempty"`
	Title         string         `json:"title,omitempty"`
	Status        int            `json:"-"`
	Detail        string         `json:"detail,omitempty"`
	Instance      string         `json:"instance,omitempty"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
	ErrorName     string         `json:"errorName,omitempty"`
	Error         error          `json:"-"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
