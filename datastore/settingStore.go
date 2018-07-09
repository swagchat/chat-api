package datastore

import "github.com/swagchat/chat-api/model"

type settingStore interface {
	createSettingStore()

	SelectLatestSetting() (*model.Setting, error)
}
