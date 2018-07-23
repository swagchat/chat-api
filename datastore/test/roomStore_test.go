package datastore_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

func TestRoomStore(t *testing.T) {
	t.Run("Insert room", func(t *testing.T) {
		r := &model.Room{}
		r.RoomID = "room-id-0001"
		r.Name = "name"
		r.MetaData = []byte(`{"key":"value"}`)
		r.Created = 123456789
		r.Modified = 123456789

		err := datastore.Provider(ctx).InsertRoom(r)
		if err != nil {
			t.Fatalf("Failed insert room test")
		}
	})
	t.Run("Select rooms", func(t *testing.T) {
		rooms, err := datastore.Provider(ctx).SelectRooms(10, 0)
		if err != nil {
			t.Fatalf("Failed Select rooms test")
		}
		if len(rooms) != 2 {
			t.Fatalf("Failed Select rooms test")
		}
	})
	t.Run("Select room", func(t *testing.T) {
		roomId := "room-id-0001"
		room, err := datastore.Provider(ctx).SelectRoom(roomId)
		if room == nil {
			t.Fatalf("Failed Select room test")
		}
		if err != nil {
			t.Fatalf("Failed Select room test")
		}
	})
	t.Run("Update room", func(t *testing.T) {
		r := &model.Room{}
		r.RoomID = "room-id-0001"
		r.Name = "name-update"
		r.MetaData = []byte(`{"key":"value"}`)
		r.Created = 123456789
		r.Modified = 123456789
		err := datastore.Provider(ctx).UpdateRoom(r)
		if err != nil {
			t.Fatalf("Failed update room test")
		}
	})
}
