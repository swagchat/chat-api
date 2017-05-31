package utils

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func CreateUuid() string {
	uuid := uuid.NewV4().String()
	return uuid
}

func CreateApiKey() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "", -1)
}
