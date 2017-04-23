package datastore

import gorp "gopkg.in/gorp.v2"

const (
	TABLE_NAME_USER         = "user"
	TABLE_NAME_ROOM         = "room"
	TABLE_NAME_ROOM_USER    = "room_user"
	TABLE_NAME_MESSAGE      = "message"
	TABLE_NAME_DEVICE       = "device"
	TABLE_NAME_SUBSCRIPTION = "subscription"
)

var dbMap *gorp.DbMap
