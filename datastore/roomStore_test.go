package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameInsertRoom       = "insert room test"
	TestNameSelectRooms      = "select rooms test"
	TestNameSelectRoom       = "select room test"
	TestNameSelectCountRooms = "select count rooms test"
	TestNameUpdateRoom       = "update room test"
)

func TestRoomStore(t *testing.T) {
	testRoomID := "room-id-0001"
	var room *model.Room
	var err error

	t.Run(TestNameInsertRoom, func(t *testing.T) {
		newRoom := &model.Room{}
		newRoom.RoomID = testRoomID
		newRoom.Name = "name"
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.Created = 123456789
		newRoom.Modified = 123456789
		newRoomUser := &model.RoomUser{}
		newRoomUser.RoomID = testRoomID
		newRoomUser.UserID = "datastore-user-id-0001"
		newRoomUser.UnreadCount = 0
		newRoomUser.Display = true
		err = Provider(ctx).InsertRoom(
			newRoom,
			InsertRoomOptionWithRoomUser([]*model.RoomUser{newRoomUser}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertRoom)
		}
	})

	t.Run(TestNameSelectRoom, func(t *testing.T) {
		room, err = Provider(ctx).SelectRoom(
			testRoomID,
			SelectRoomOptionWithUsers(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoom)
		}
		if room == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoom)
		}
	})

	t.Run(TestNameUpdateRoom, func(t *testing.T) {
		room.Name = "name-update"
		err = Provider(ctx).UpdateRoom(room)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
		updatedRoom, err := Provider(ctx).SelectRoom(testRoomID)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
		if updatedRoom == nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
		if updatedRoom.Name != "name-update" {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}

		updatedRoom.Deleted = 1
		err = Provider(ctx).UpdateRoom(updatedRoom)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
		deletedRoom, err := Provider(ctx).SelectRoom(testRoomID)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
		if deletedRoom != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoom)
		}
	})

	t.Run(TestNameSelectRooms, func(t *testing.T) {
		rooms, err := Provider(ctx).SelectRooms(
			0,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if len(rooms) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}

		rooms, err = Provider(ctx).SelectRooms(
			10,
			20,
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if len(rooms) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}

		orderInfo1 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		orders := []*scpb.OrderInfo{orderInfo1}
		rooms, err = Provider(ctx).SelectRooms(
			10,
			0,
			SelectRoomsOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if len(rooms) != 10 {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if rooms[0].RoomID != "datastore-room-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}

		orderInfo2 := &scpb.OrderInfo{
			Field: "last_message_updated",
			Order: scpb.Order_Desc,
		}
		orderInfo3 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		orders = []*scpb.OrderInfo{orderInfo2, orderInfo3}
		rooms, err = Provider(ctx).SelectRooms(
			20,
			0,
			SelectRoomsOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if len(rooms) != 20 {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if rooms[0].RoomID != "datastore-room-id-0010" {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if rooms[9].RoomID != "datastore-room-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if rooms[10].RoomID != "datastore-room-id-0011" {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
		if rooms[19].RoomID != "datastore-room-id-0020" {
			t.Fatalf("Failed to %s", TestNameSelectRooms)
		}
	})

	t.Run(TestNameSelectCountRooms, func(t *testing.T) {
		count, err := Provider(ctx).SelectCountRooms()
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectCountRooms)
		}
		if count != 20 {
			t.Fatalf("Failed to %s", TestNameSelectCountRooms)
		}
	})
}
