package datastore

import (
	"context"
	"fmt"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
	"gopkg.in/gorp.v2"
)

func rdbCreateRoomStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateRoomStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.Room{}, tableNameRoom)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "room_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating room table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}
}

func rdbInsertRoom(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, room *model.Room, opts ...InsertRoomOption) error {
	span := tracer.StartSpan(ctx, "rdbInsertRoom", "datastore")
	defer tracer.Finish(span)

	opt := insertRoomOptions{}
	for _, o := range opts {
		o(&opt)
	}

	err := tx.Insert(room)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting room")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	for _, ru := range opt.users {
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id=?;", tableNameRoomUser)
		_, err = tx.Exec(query, room.RoomID)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while inserting room")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}

		err = tx.Insert(ru)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while inserting room")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
	}

	return nil
}

func rdbSelectRooms(ctx context.Context, dbMap *gorp.DbMap, limit, offset int32, opts ...SelectRoomsOption) ([]*model.Room, error) {
	span := tracer.StartSpan(ctx, "rdbSelectRooms", "datastore")
	defer tracer.Finish(span)

	opt := selectRoomsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var rooms []*model.Room
	query := fmt.Sprintf(`SELECT
	room_id,
	user_id,
	name,
	picture_url,
	information_url,
	meta_data,
	type,
	last_message,
	last_message_updated,
	created,
	modified
	FROM %s
	WHERE deleted = 0`, tableNameRoom)
	params := make(map[string]interface{})

	query = fmt.Sprintf("%s ORDER BY", query)
	if opt.orders == nil {
		query = fmt.Sprintf("%s created DESC", query)
	} else {
		i := 1
		for _, orderInfo := range opt.orders {
			query = fmt.Sprintf("%s %s %s", query, orderInfo.Field, orderInfo.Order.String())
			if i < len(opt.orders) {
				query = fmt.Sprintf("%s,", query)
			}
			i++
		}
	}

	query = fmt.Sprintf("%s LIMIT :limit OFFSET :offset", query)
	params["limit"] = limit
	params["offset"] = offset
	_, err := dbMap.Select(&rooms, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting rooms")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return rooms, nil
}

func rdbSelectRoom(ctx context.Context, dbMap *gorp.DbMap, roomID string, opts ...SelectRoomOption) (*model.Room, error) {
	span := tracer.StartSpan(ctx, "rdbSelectRoom", "datastore")
	defer tracer.Finish(span)

	opt := selectRoomOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var rooms []*model.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted=0;", tableNameRoom)
	params := map[string]interface{}{"roomId": roomID}
	_, err := dbMap.Select(&rooms, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting room")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	var room *model.Room
	if len(rooms) != 1 {
		return nil, nil
	}

	room = rooms[0]

	if opt.withUsers {
		users, err := rdbSelectUsersForRoom(ctx, dbMap, roomID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting room")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return nil, err
		}
		room.Users = users
	}

	return room, nil
}

func rdbSelectUsersForRoom(ctx context.Context, dbMap *gorp.DbMap, roomID string) ([]*model.MiniUser, error) {
	span := tracer.StartSpan(ctx, "rdbSelectUsersForRoom", "datastore")
	defer tracer.Finish(span)

	var users []*model.MiniUser
	query := fmt.Sprintf(`SELECT
u.user_id,
u.name,
u.picture_url,
u.information_url,
u.meta_data,
u.can_block,
u.last_accessed,
u.created,
u.modified,
ru.display AS ru_display
FROM %s AS ru 
LEFT JOIN %s AS u ON ru.user_id = u.user_id 
WHERE ru.room_id = :roomId AND u.deleted = 0 
ORDER BY u.created;`, tableNameRoomUser, tableNameUser)
	params := map[string]interface{}{"roomId": roomID}
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting users for room")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return users, nil
}

func rdbSelectCountRooms(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectRoomsOption) (int64, error) {
	span := tracer.StartSpan(ctx, "rdbSelectCountRooms", "datastore")
	defer tracer.Finish(span)

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE deleted = 0;", tableNameRoom)
	count, err := dbMap.SelectInt(query)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting room count")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return 0, err
	}

	return count, nil
}

func rdbUpdateRoom(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, room *model.Room, opts ...UpdateRoomOption) error {
	span := tracer.StartSpan(ctx, "rdbUpdateRoom", "datastore")
	defer tracer.Finish(span)

	opt := updateRoomOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if room.DeletedTimestamp != 0 {
		return rdbUpdateRoomDeleted(ctx, dbMap, tx, room)
	}

	_, err := tx.Update(room)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating room")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	for _, ru := range opt.users {
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id=?;", tableNameRoomUser)
		_, err := tx.Exec(query, room.RoomID)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while inserting room")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}

		err = tx.Insert(ru)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while updating room")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
	}

	return nil
}

func rdbUpdateRoomDeleted(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, room *model.Room) error {
	span := tracer.StartSpan(ctx, "rdbUpdateRoomDeleted", "datastore")
	defer tracer.Finish(span)

	err := rdbDeleteRoomUsers(
		ctx,
		dbMap,
		tx,
		DeleteRoomUsersOptionFilterByRoomIDs([]string{room.RoomID}),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(room.DeletedTimestamp),
		DeleteSubscriptionsOptionFilterByRoomID(room.RoomID),
	)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET deleted=? WHERE room_id=?;", tableNameRoom)
	_, err = tx.Exec(query, room.DeletedTimestamp, room.RoomID)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting room")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}
