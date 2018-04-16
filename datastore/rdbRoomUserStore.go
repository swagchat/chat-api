package datastore

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateRoomUserStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.RoomUser{}, tableNameRoomUser)
	tableMap.SetUniqueTogether("room_id", "user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create room user table error",
			Error:   err,
		})
	}
}

func rdbDeleteAndInsertRoomUsers(db string, roomUsers []*models.RoomUser) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("DELETE FROM ", tableNameRoomUser, " WHERE room_id=:roomId;")
	params := map[string]interface{}{"roomId": roomUsers[0].RoomID}
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

func rdbInsertRoomUsers(db string, roomUsers []*models.RoomUser) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, roomUser := range roomUsers {
		ru, err := rdbSelectRoomUser(db, roomUser.RoomID, roomUser.UserID)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while creating room's users")
		}
		if ru == nil {
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

func rdbSelectRoomUser(db, roomID, userID string) (*models.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", tableNameRoomUser, " WHERE room_id=:roomId AND user_id=:userId;")
	params := map[string]interface{}{
		"roomId": roomID,
		"userId": userID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectRoomUserOfOneOnOne(db, myUserID, opponentUserID string) (*models.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", tableNameRoomUser, " WHERE room_id IN (SELECT room_id FROM ", tableNameRoom, " WHERE type=:type AND user_id=:myUserId) AND user_id=:opponentUserId;")
	params := map[string]interface{}{
		"type":           models.OneOnOne,
		"myUserId":       myUserID,
		"opponentUserId": opponentUserID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's user")
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectRoomUsersByRoomID(db, roomID string) ([]*models.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT room_id, user_id, unread_count, meta_data, created, modified FROM ", tableNameRoomUser, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.New("An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func rdbSelectRoomUsersByUserID(db, userID string) ([]*models.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", tableNameRoomUser, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.New("An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func rdbSelectRoomUsersByRoomIDAndUserIDs(db string, roomID *string, userIDs []string) ([]*models.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*models.RoomUser
	var userIDsQuery string
	var userIDsParams map[string]interface{}
	var roomIDParams map[string]interface{}
	if userIDs != nil {
		userIDsQuery, userIDsParams = utils.MakePrepareForInExpression(userIDs)
	}
	if roomID != nil {
		roomIDParams = map[string]interface{}{"roomId": roomID}
	}
	params := utils.MergeMap(userIDsParams, roomIDParams)

	query := utils.AppendStrings("SELECT * ",
		"FROM ", tableNameRoomUser,
		" WHERE ")
	if roomID != nil {
		query = utils.AppendStrings(query, " room_id=:roomId")
	}
	if roomID != nil && userIDs != nil {
		query = utils.AppendStrings(query, " AND ")
	}
	if userIDs != nil {
		query = utils.AppendStrings(query, " user_id IN (", userIDsQuery, ")")
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting room's users")
	}

	return roomUsers, nil
}

func rdbUpdateRoomUser(db string, roomUser *models.RoomUser) (*models.RoomUser, error) {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
	}

	updateQuery := ""
	params := map[string]interface{}{
		"roomId": roomUser.RoomID,
		"userId": roomUser.UserID,
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
		query := utils.AppendStrings("UPDATE ", tableNameRoomUser, " SET "+updateQuery+" WHERE room_id=:roomId AND user_id=:userId;")
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return nil, errors.Wrap(err, "An error occurred while updating room's users")
		}

		if roomUser.UnreadCount != nil {
			query = utils.AppendStrings("UPDATE ", tableNameUser,
				" SET unread_count=(SELECT SUM(unread_count) FROM ", tableNameRoomUser,
				" WHERE user_id=:userId1) WHERE user_id=:userId2;")
			params = map[string]interface{}{
				"userId1": roomUser.UserID,
				"userId2": roomUser.UserID,
			}
			_, err = trans.Exec(query, params)
			if err != nil {
				err = trans.Rollback()
				return nil, errors.Wrap(err, "An error occurred while updating user unread count")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return nil, errors.New("An error occurred while commit updating room's user")
	}

	return roomUser, nil
}

func rdbDeleteRoomUser(db, roomID string, userIDs []string) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	var query string
	var params map[string]interface{}
	if userIDs == nil {
		query = utils.AppendStrings("DELETE FROM ", tableNameRoomUser, " WHERE room_id=:roomId;")
		params = map[string]interface{}{"roomId": roomID}
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while deleting room's users")
		}

		query = utils.AppendStrings("UPDATE ", tableNameSubscription, " SET deleted=:deleted WHERE room_id=:roomId;")
		params = map[string]interface{}{
			"roomId":  roomID,
			"deleted": time.Now().Unix(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while updating subsctiptions")
		}
	} else {
		var userIdsQuery string
		userIdsQuery, params = utils.MakePrepareForInExpression(userIDs)
		query = utils.AppendStrings("DELETE FROM ", tableNameRoomUser, " WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		params["roomId"] = roomID
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while deleting room's users")
		}

		query = utils.AppendStrings("UPDATE ", tableNameSubscription, " SET deleted=:deleted WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
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
