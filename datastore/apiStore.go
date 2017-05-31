package datastore

type ApiStore interface {
	CreateApiStore()

	InsertApi(name string) StoreResult
	SelectLatestApi(name string) StoreResult
}
