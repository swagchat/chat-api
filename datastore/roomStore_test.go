package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

func TestRoomStore(t *testing.T) {
	testRoomID := "room-id-0001"

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
	err := Provider(ctx).InsertRoom(
		newRoom,
		InsertRoomOptionWithRoomUser([]*model.RoomUser{newRoomUser}),
	)
	if err != nil {
		t.Fatalf("Failed insert room test")
	}

	room, err := Provider(ctx).SelectRoom(
		testRoomID,
		SelectRoomOptionWithUsers(true),
	)
	if err != nil {
		t.Fatalf("Failed select room test")
	}
	if room == nil {
		t.Fatalf("Failed select room test")
	}

	room.Name = "name-update"
	err = Provider(ctx).UpdateRoom(room)
	if err != nil {
		t.Fatalf("Failed update room test")
	}
	updatedRoom, err := Provider(ctx).SelectRoom(testRoomID)
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
	err = Provider(ctx).UpdateRoom(updatedRoom)
	if err != nil {
		t.Fatalf("Failed update room test")
	}
	deletedRoom, err := Provider(ctx).SelectRoom(testRoomID)
	if err != nil {
		t.Fatalf("Failed select room test")
	}
	if deletedRoom != nil {
		t.Fatalf("Failed delete room test")
	}

	rooms, err := Provider(ctx).SelectRooms(
		0,
		0,
	)
	if err != nil {
		t.Fatalf("Failed select rooms test")
	}
	if len(rooms) != 0 {
		t.Fatalf("Failed select rooms test")
	}

	rooms, err = Provider(ctx).SelectRooms(
		10,
		20,
	)
	if err != nil {
		t.Fatalf("Failed select rooms test")
	}
	if len(rooms) != 0 {
		t.Fatalf("Failed select rooms test")
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
		t.Fatalf("Failed select rooms test")
	}
	if len(rooms) != 10 {
		t.Fatalf("Failed select rooms test")
	}
	if rooms[0].RoomID != "datastore-room-id-0001" {
		t.Fatalf("Failed select rooms test")
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
		t.Fatalf("Failed select rooms test")
	}
	if len(rooms) != 20 {
		t.Fatalf("Failed select rooms test")
	}
	if rooms[0].RoomID != "datastore-room-id-0010" {
		t.Fatalf("Failed select rooms test")
	}
	if rooms[9].RoomID != "datastore-room-id-0001" {
		t.Fatalf("Failed select rooms test")
	}
	if rooms[10].RoomID != "datastore-room-id-0011" {
		t.Fatalf("Failed select rooms test")
	}
	if rooms[19].RoomID != "datastore-room-id-0020" {
		t.Fatalf("Failed select rooms test")
	}

	count, err := Provider(ctx).SelectCountRooms()
	if err != nil {
		t.Fatalf("Failed select count rooms test")
	}
	if count != 20 {
		t.Fatalf("Failed select count rooms test")
	}
}
