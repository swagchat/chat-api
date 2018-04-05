package datastore

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
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
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create room table error",
			Error:   err,
		})
	}
}

func RdbInsertRoom(room *models.Room) (*models.Room, error) {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
	}

	err = trans.Insert(room)
	if err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while creating room")
	}

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
	for _, userID := range room.RequestRoomUserIds.UserIds {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomId:      room.RoomId,
			UserId:      userID,
			UnreadCount: &zero,
			MetaData:    []byte("{}"),
			Created:     room.Created,
			Modified:    room.Modified,
		})
	}

	for _, roomUser := range roomUsers {
		err = trans.Insert(roomUser)
		if err != nil {
			err = trans.Rollback()
			return nil, errors.Wrap(err, "An error occurred while creating room's users")
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while commit creating room")
	}

	return room, nil
}

func RdbSelectRoom(roomId string) (*models.Room, error) {
	slave := RdbStoreInstance().replica()

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": roomId}
	_, err := slave.Select(&rooms, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room")
	}

	if len(rooms) == 1 {
		return rooms[0], nil
	}

	return nil, nil
}

func RdbSelectRooms() ([]*models.Room, error) {
	slave := RdbStoreInstance().replica()

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT room_id, user_id, name, picture_url, information_url, meta_data, type, last_message, last_message_updated, created, modified FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	_, err := slave.Select(&rooms, query)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting rooms")
	}

	return rooms, nil
}

func RdbSelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.UserForRoom
	query := utils.AppendStrings("SELECT ",
		"u.user_id, ",
		"u.name, ",
		"u.picture_url, ",
		"u.information_url, ",
		"u.meta_data, ",
		"u.is_bot, ",
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
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	return users, nil
}

func RdbSelectCountRooms() (int64, error) {
	slave := RdbStoreInstance().replica()

	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", TABLE_NAME_ROOM, " WHERE deleted = 0;")
	count, err := slave.SelectInt(query)
	if err != nil {
		return 0, errors.Wrap(err, "An error occurred while getting room count")
	}

	return count, nil
}

func RdbUpdateRoom(room *models.Room) (*models.Room, error) {
	master := RdbStoreInstance().master()

	_, err := master.Update(room)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while updating room")
	}

	return room, nil
}

func RdbUpdateRoomDeleted(roomId string) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating room")
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating room")
	}

	err = trans.Commit()
	if err != nil {
		err := trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit updating room ")
	}

	return nil
}
