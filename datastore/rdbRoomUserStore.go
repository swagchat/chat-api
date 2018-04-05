package datastore

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateRoomUserStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.RoomUser{}, TABLE_NAME_ROOM_USER)
	tableMap.SetUniqueTogether("room_id", "user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create room user table error",
			Error:   err,
		})
	}
}

func RdbDeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{"roomId": roomUsers[0].RoomId}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting room's users")
	}

	for _, roomUser := range roomUsers {
		err = trans.Insert(roomUser)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while creating room's users")
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating room's users")
	}

	return nil
}

func RdbInsertRoomUsers(roomUsers []*models.RoomUser) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, roomUser := range roomUsers {
		roomUser, err := RdbSelectRoomUser(roomUser.RoomId, roomUser.UserId)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while creating room's users")
		}
		if roomUser == nil {
			err = trans.Insert(roomUser)
			if err != nil {
				err = trans.Rollback()
				return errors.Wrap(err, "An error occurred while creating room's user")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating room's user items")
	}

	return nil
}

func RdbSelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	slave := RdbStoreInstance().replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id=:userId;")
	params := map[string]interface{}{
		"roomId": roomId,
		"userId": userId,
	}
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	slave := RdbStoreInstance().replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id IN (SELECT room_id FROM ", TABLE_NAME_ROOM, " WHERE type=:type AND user_id=:myUserId) AND user_id=:opponentUserId;")
	params := map[string]interface{}{
		"type":           models.ONE_ON_ONE,
		"myUserId":       myUserId,
		"opponentUserId": opponentUserId,
	}
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's user")
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func RdbSelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	slave := RdbStoreInstance().replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT room_id, user_id, unread_count, meta_data, created, modified FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.New("An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func RdbSelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	slave := RdbStoreInstance().replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.New("An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func RdbSelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	slave := RdbStoreInstance().replica()

	var roomUsers []*models.RoomUser
	var userIdsQuery string
	var userIdsParams map[string]interface{}
	var roomIdParams map[string]interface{}
	if userIds != nil {
		userIdsQuery, userIdsParams = utils.MakePrepareForInExpression(userIds)
	}
	if roomId != nil {
		roomIdParams = map[string]interface{}{"roomId": roomId}
	}
	params := utils.MergeMap(userIdsParams, roomIdParams)

	query := utils.AppendStrings("SELECT * ",
		"FROM ", TABLE_NAME_ROOM_USER,
		" WHERE ")
	if roomId != nil {
		query = utils.AppendStrings(query, " room_id=:roomId")
	}
	if roomId != nil && userIds != nil {
		query = utils.AppendStrings(query, " AND ")
	}
	if userIds != nil {
		query = utils.AppendStrings(query, " user_id IN (", userIdsQuery, ")")
	}
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func RdbUpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
	}

	updateQuery := ""
	params := map[string]interface{}{
		"roomId": roomUser.RoomId,
		"userId": roomUser.UserId,
	}
	if roomUser.UnreadCount != nil {
		params["unreadCount"] = roomUser.UnreadCount
		updateQuery = "unread_count=:unreadCount"
	}
	if roomUser.MetaData != nil {
		params["metaData"] = roomUser.MetaData
		if updateQuery == "" {
			updateQuery = "meta_data=:metaData"
		} else {
			updateQuery = utils.AppendStrings(updateQuery, ",", "meta_data=:metaData")
		}
	}
	if updateQuery != "" {
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET "+updateQuery+" WHERE room_id=:roomId AND user_id=:userId;")
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return nil, errors.Wrap(err, "An error occurred while updating room's users")
		}

		if roomUser.UnreadCount != nil {
			query = utils.AppendStrings("UPDATE ", TABLE_NAME_USER,
				" SET unread_count=(SELECT SUM(unread_count) FROM ", TABLE_NAME_ROOM_USER,
				" WHERE user_id=:userId1) WHERE user_id=:userId2;")
			params = map[string]interface{}{
				"userId1": roomUser.UserId,
				"userId2": roomUser.UserId,
			}
			_, err = trans.Exec(query, params)
			if err != nil {
				err = trans.Rollback()
				return nil, errors.Wrap(err, "An error occurred while updating user unread count")
			}
		}
	}

	if roomUser == nil {
		err = trans.Commit()
		if err != nil {
			err = trans.Rollback()
			return nil, errors.New("An error occurred while commit updating room's user")
		}
	}

	return roomUser, nil
}

func RdbDeleteRoomUser(roomId string, userIds []string) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	var query string
	var params map[string]interface{}
	if userIds == nil {
		query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
		params = map[string]interface{}{"roomId": roomId}
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while deleting room's users")
		}

		query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId;")
		params = map[string]interface{}{
			"roomId":  roomId,
			"deleted": time.Now().Unix(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while updating subsctiptions")
		}
	} else {
		var userIdsQuery string
		userIdsQuery, params = utils.MakePrepareForInExpression(userIds)
		query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		params["roomId"] = roomId
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while deleting room's users")
		}

		query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		params["deleted"] = time.Now().Unix()
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while updating subsctiptions")
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.New("An error occurred while commit deleting room's users")
	}

	return nil
}
