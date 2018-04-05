package datastore

import "github.com/swagchat/chat-api/models"

type SettingStore interface {
	CreateSettingStore()

	SelectLatestSetting() (*models.Setting, error)
}
