package datastore

import (
	"fmt"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
	"github.com/pkg/errors"
)

func rdbCreateGuestSettingStore(db string) {
	master := RdbStore(db).master()
	tableMap := master.AddTableWithName(model.GuestSetting{}, tableNameGuestSetting)
	tableMap.SetKeys(true, "id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating guest setting table. %v.", err))
		return
	}
}

func rdbInsertGuestSetting(db string, setting *model.GuestSetting) (*model.GuestSetting, error) {
	master := RdbStore(db).master()

	if err := master.Insert(setting); err != nil {
		err = errors.Wrap(err, "An error occurred while inserting guest setting")
		logger.Error(err.Error())
		return nil, err
	}

	return setting, nil
}

func rdbSelectGuestSetting(db string) (*model.GuestSetting, error) {
	slave := RdbStore(db).replica()

	var ss []*model.GuestSetting
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=:id;", tableNameGuestSetting)
	params := map[string]interface{}{"id": 1}
	_, err := slave.Select(&ss, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting guest setting")
		logger.Error(err.Error())
		return nil, err
	}

	if len(ss) == 1 {
		return ss[0], nil
	}

	return nil, nil
}

func rdbUpdateGuestSetting(db string, setting *model.GuestSetting) error {
	master := RdbStore(db).master()

	setting.ID = 1
	if _, err := master.Update(setting); err != nil {
		err = errors.Wrap(err, "An error occurred while updating guest setting")
		logger.Error(err.Error())
		return err
	}

	return nil
}
