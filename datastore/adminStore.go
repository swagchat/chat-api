package datastore

type AdminStore interface {
	CreateAdminStore()

	InsertAdmin() StoreResult
	SelectLatestAdmin() StoreResult
}
