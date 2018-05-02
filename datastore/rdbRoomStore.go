package datastore

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateRoomStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Room{}, tableNameRoom)
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

func rdbInsertRoom(db string, room *models.Room) (*models.Room, error) {
	master := RdbStore(db).master()
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
		RoomID:      room.RoomID,
		UserID:      room.UserID,
		UnreadCount: &zero,
		MetaData:    []byte("{}"),
		Created:     room.Created,
		Modified:    room.Modified,
	}
	roomUsers = append(roomUsers, roomUser)
	for _, userID := range room.RequestRoomUserIDs.UserIDs {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomID:      room.RoomID,
			UserID:      userID,
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

func rdbSelectRoom(db, roomID string) (*models.Room, error) {
	replica := RdbStore(db).replica()

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", tableNameRoom, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": roomID}
	_, err := replica.Select(&rooms, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room")
	}

	if len(rooms) == 1 {
		return rooms[0], nil
	}

	return nil, nil
}

func rdbSelectRooms(db string) ([]*models.Room, error) {
	replica := RdbStore(db).replica()

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT room_id, user_id, name, picture_url, information_url, meta_data, type, last_message, last_message_updated, created, modified FROM ", tableNameRoom, " WHERE deleted = 0;")
	_, err := replica.Select(&rooms, query)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting rooms")
	}

	return rooms, nil
}

func rdbSelectUsersForRoom(db, roomID string) ([]*models.UserForRoom, error) {
	replica := RdbStore(db).replica()

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
		"u.last_accessed, ",
		"u.created, ",
		"u.modified, ",
		"ru.unread_count AS ru_unread_count, ",
		"ru.meta_data AS ru_meta_data, ",
		"ru.created AS ru_created, ",
		"ru.modified AS ru_modified ",
		"FROM ", tableNameRoomUser, " AS ru ",
		"LEFT JOIN ", tableNameUser, " AS u ",
		"ON ru.user_id = u.user_id ",
		"WHERE ru.room_id = :roomId AND u.deleted = 0 ",
		"ORDER BY u.created;")
	params := map[string]interface{}{"roomId": roomID}
	_, err := replica.Select(&users, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	return users, nil
}

func rdbSelectCountRooms(db string) (int64, error) {
	replica := RdbStore(db).replica()

	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", tableNameRoom, " WHERE deleted = 0;")
	count, err := replica.SelectInt(query)
	if err != nil {
		return 0, errors.Wrap(err, "An error occurred while getting room count")
	}

	return count, nil
}

func rdbUpdateRoom(db string, room *models.Room) (*models.Room, error) {
	master := RdbStore(db).master()

	_, err := master.Update(room)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while updating room")
	}

	return room, nil
}

func rdbUpdateRoomDeleted(db, roomID string) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("DELETE FROM ", tableNameRoomUser, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating room")
	}

	query = utils.AppendStrings("UPDATE ", tableNameSubscription, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomID,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = utils.AppendStrings("UPDATE ", tableNameRoom, " SET deleted=:deleted WHERE room_id=:roomId;")
	params = map[string]interface{}{
		"roomId":  roomID,
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
