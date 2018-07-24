package datastore_test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

func TestRoomStore(t *testing.T) {
	newRoom := &model.Room{}
	newRoom.RoomID = "room-id-0001"
	newRoom.Name = "name"
	newRoom.MetaData = []byte(`{"key":"value"}`)
	newRoom.Created = 123456789
	newRoom.Modified = 123456789
	err := datastore.Provider(ctx).InsertRoom(newRoom)
	if err != nil {
		t.Fatalf("Failed insert room test")
	}

	rooms, err := datastore.Provider(ctx).SelectRooms(10, 0)
	if err != nil {
		t.Fatalf("Failed select rooms test")
	}
	if len(rooms) != 2 {
		t.Fatalf("Failed select rooms test")
	}

	room, err := datastore.Provider(ctx).SelectRoom("room-id-0001")
	if err != nil {
		t.Fatalf("Failed select room test")
	}
	if room == nil {
		t.Fatalf("Failed select room test")
	}

	room.Name = "name-update"
	err = datastore.Provider(ctx).UpdateRoom(room)
	if err != nil {
		t.Fatalf("Failed update room test")
	}
	updatedRoom, err := datastore.Provider(ctx).SelectRoom("room-id-0001")
	if err != nil {
		t.Fatalf("Failed select room test")
	}
	if updatedRoom == nil {
		t.Fatalf("Failed select room test")
	}
	if updatedRoom.Name != "name-update" {
		t.Fatalf("Failed update room test")
	}

	updatedRoom.Deleted = 1
	err = datastore.Provider(ctx).UpdateRoom(updatedRoom)
	if err != nil {
		t.Fatalf("Failed update room test")
	}
	deletedRoom, err := datastore.Provider(ctx).SelectRoom("room-id-0001")
	if err != nil {
		t.Fatalf("Failed select room test")
	}
	if deletedRoom != nil {
		t.Fatalf("Failed delete room test")
	}
}

func TestRoomsStore(t *testing.T) {
	var newRoom *model.Room
	nowTimestamp := time.Now().Unix()
	for i := 1; i <= 10; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("rooms-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp + int64(i)
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			t.Fatalf("Failed insert room test")
		}
	}
	for i := 11; i <= 20; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("rooms-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			t.Fatalf("Failed insert room test")
		}
	}

	orders := make(map[string]scpb.Order, 1)
	rooms, err := datastore.Provider(ctx).SelectRooms(
		0,
		0,
	)
	if err != nil {
		t.Fatalf("Failed Select rooms test")
	}
	if len(rooms) != 0 {
		t.Fatalf("Failed Select rooms test")
	}

	orders = make(map[string]scpb.Order, 1)
	orders["created"] = scpb.Order_Asc
	rooms, err = datastore.Provider(ctx).SelectRooms(
		10,
		0,
		datastore.RoomOptionOrders(orders),
	)
	if err != nil {
		t.Fatalf("Failed Select rooms test")
	}
	if len(rooms) != 10 {
		t.Fatalf("Failed Select rooms test")
	}
	if rooms[0].RoomID != "room-id-0000" {
		t.Fatalf("Failed Select rooms test")
	}

	orders = make(map[string]scpb.Order, 1)
	orders["last_message_updated"] = scpb.Order_Desc
	orders["created"] = scpb.Order_Asc
	rooms, err = datastore.Provider(ctx).SelectRooms(
		20,
		0,
		datastore.RoomOptionOrders(orders),
	)
	if err != nil {
		t.Fatalf("Failed Select rooms test")
	}
	if len(rooms) != 20 {
		t.Fatalf("Failed Select rooms test")
	}
	if rooms[0].RoomID != "rooms-id-0010" {
		t.Fatalf("Failed Select rooms test")
	}
	if rooms[11].RoomID != "rooms-id-0011" {
		t.Fatalf("Failed Select rooms test")
	}
	if rooms[19].RoomID != "rooms-id-0019" {
		t.Fatalf("Failed Select rooms test")
	}
}
