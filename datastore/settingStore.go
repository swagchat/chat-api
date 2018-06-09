package datastore

import "github.com/swagchat/chat-api/models"

type settingStore interface {
	createSettingStore()

	SelectLatestSetting() (*models.Setting, error)
}
