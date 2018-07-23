package utils

import (
	uuid "github.com/satori/go.uuid"
)

// GenerateUUID is generate UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	return uuid
}
