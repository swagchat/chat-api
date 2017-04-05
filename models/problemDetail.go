package models

const (
	ERROR_NAME_INVALID_JSON       = "invalid-json"
	ERROR_NAME_INVALID_PARAM      = "invalid-param"
	ERROR_NAME_DATABASE_ERROR     = "database-error"
	ERROR_NAME_NOTIFICATION_ERROR = "notification-error"
)

type ProblemDetail struct {
	Type          string         `json:"type,omitempty"`
	Title         string         `json:"title"`
	Status        int            `json:"status,omitempty"`
	Detail        string         `json:"detail,omitempty"`
	Instance      string         `json:"instance,omitempty"`
	ErrorName     string         `json:"errorName"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
