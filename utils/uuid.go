package utils

import uuid "github.com/satori/go.uuid"

func CreateUuid() string {
	uuid := uuid.NewV4().String()
	return uuid
}
