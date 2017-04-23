package datastore

import (
	"log"
	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateRoomStore() {
	tableMap := dbMap.AddTableWithName(models.Room{}, TABLE_NAME_ROOM)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "room_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertRoom(room *models.Room) StoreResult {
	result := StoreResult{}
	if err := dbMap.Insert(room); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating room item.", err)
	}
	result.Data = room
	return result
}

func RdbSelectRoom(roomId string) StoreResult {
	result := StoreResult{}
	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": roomId}
	if _, err := dbMap.Select(&rooms, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room item.", err)
	}
	if len(rooms) == 1 {
		result.Data = rooms[0]
	}
	return result
}

func RdbSelectRooms() StoreResult {
	result := StoreResult{}
	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	_, err := dbMap.Select(&rooms, query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room items.", err)
	}
	result.Data = rooms
	return result
}

func RdbSelectUsersForRoom(roomId string) StoreResult {
	result := StoreResult{}
	var users []*models.UserForRoom
	query := utils.AppendStrings("SELECT ",
		"u.user_id, ",
		"u.name, ",
		"u.picture_url, ",
		"u.information_url, ",
		"u.meta_data, ",
		"u.created, ",
		"u.modified, ",
		"ru.unread_count AS ru_unread_count, ",
		"ru.meta_data AS ru_meta_data, ",
		"ru.created AS ru_created ",
		"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
		"LEFT JOIN ", TABLE_NAME_USER, " AS u ",
		"ON ru.user_id = u.user_id ",
		"WHERE ru.room_id = :roomId AND u.deleted = 0 ",
		"ORDER BY u.created;")
	params := map[string]interface{}{"roomId": roomId}
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room users.", err)
	}
	result.Data = users
	return result
}

func RdbSelectCountRooms() StoreResult {
	result := StoreResult{}
	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	count, err := dbMap.SelectInt(query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room count.", err)
	}
	result.Data = count
	return result
}

func RdbUpdateRoom(room *models.Room) StoreResult {
	result := StoreResult{}
	_, err := dbMap.Update(room)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
	}
	result.Data = room
	return result
}

func RdbUpdateRoomDeleted(roomId string) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating room's user items.", err)
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().UnixNano(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().UnixNano(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room item.", err)
		}
	} else {
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room item.", err)
		}
	}
	return result
}
