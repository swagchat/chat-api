package datastore

type SettingStore interface {
	CreateSettingStore()

	SelectLatestSetting() StoreResult
}
