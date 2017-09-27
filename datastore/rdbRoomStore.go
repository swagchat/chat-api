package datastore

import (
	"log"
	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateRoomStore() {
	master := RdbStoreInstance().master()
	tableMap := master.AddTableWithName(models.Room{}, TABLE_NAME_ROOM)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "room_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertRoom(room *models.Room) StoreResult {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	result := StoreResult{}
	if err = master.Insert(room); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating room item.", err)
		if err = trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room item.", err)
		}
		return result
	}
	result.Data = room

	var zero int64
	zero = 0
	roomUsers := make([]*models.RoomUser, 0)
	roomUser := &models.RoomUser{
		RoomId:      room.RoomId,
		UserId:      room.UserId,
		UnreadCount: &zero,
		MetaData:    []byte("{}"),
		Created:     room.Created,
		Modified:    room.Modified,
	}
	roomUsers = append(roomUsers, roomUser)
	for _, userId := range room.RequestRoomUserIds.UserIds {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomId:      room.RoomId,
			UserId:      userId,
			UnreadCount: &zero,
			MetaData:    []byte("{}"),
			Created:     room.Created,
			Modified:    room.Modified,
		})
	}

	for _, roomUser := range roomUsers {
		if err := trans.Insert(roomUser); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room's user items.", err)
			}
			return result
		}
	}

	if result.ProblemDetail == nil {
		if err = trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit creating room item.", err)
		}
	}
	return result
}

func RdbSelectRoom(roomId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": roomId}
	if _, err := slave.Select(&rooms, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room item.", err)
	}
	if len(rooms) == 1 {
		result.Data = rooms[0]
	}
	return result
}

func RdbSelectRooms() StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var rooms []*models.Room
	query := utils.AppendStrings("SELECT room_id, user_id, name, picture_url, information_url, meta_data, type, last_message, last_message_updated, created, modified FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	_, err := slave.Select(&rooms, query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room items.", err)
	}
	result.Data = rooms
	return result
}

func RdbSelectUsersForRoom(roomId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var users []*models.UserForRoom
	query := utils.AppendStrings("SELECT ",
		"u.user_id, ",
		"u.name, ",
		"u.picture_url, ",
		"u.information_url, ",
		"u.meta_data, ",
		"u.is_can_block, ",
		"u.is_show_users, ",
		"u.created, ",
		"u.modified, ",
		"ru.unread_count AS ru_unread_count, ",
		"ru.meta_data AS ru_meta_data, ",
		"ru.created AS ru_created, ",
		"ru.modified AS ru_modified ",
		"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
		"LEFT JOIN ", TABLE_NAME_USER, " AS u ",
		"ON ru.user_id = u.user_id ",
		"WHERE ru.room_id = :roomId AND u.deleted = 0 ",
		"ORDER BY u.created;")
	params := map[string]interface{}{"roomId": roomId}
	_, err := slave.Select(&users, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's users.", err)
	}
	result.Data = users
	return result
}

func RdbSelectCountRooms() StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	count, err := slave.SelectInt(query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room count.", err)
	}
	result.Data = count
	return result
}

func RdbUpdateRoom(room *models.Room) StoreResult {
	master := RdbStoreInstance().master()
	result := StoreResult{}
	_, err := master.Update(room)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating room item.", err)
	}
	result.Data = room
	return result
}

func RdbUpdateRoomDeleted(roomId string) StoreResult {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating room item.", err)
		}
		return result
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating room item.", err)
		}
		return result
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating room item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating room item.", err)
		}
		return result
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit updating room item.", err)
		}
	}
	return result
}
