package datastore

import (
	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createMessageStore() {
	master := RdbStore(p.database).master()
	rdbCreateMessageStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertMessage(message *model.Message) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertMessage(p.ctx, master, tx, message)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *gcpSQLProvider) SelectMessages(limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectMessages(p.ctx, replica, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectMessage(messageID string) (*model.Message, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectMessage(p.ctx, replica, messageID)
}

func (p *gcpSQLProvider) SelectCountMessages(opts ...SelectMessagesOption) (int64, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectCountMessages(p.ctx, replica, opts...)
}

func (p *gcpSQLProvider) UpdateMessage(message *model.Message) error {
	master := RdbStore(p.database).master()
	return rdbUpdateMessage(p.ctx, master, message)
}
