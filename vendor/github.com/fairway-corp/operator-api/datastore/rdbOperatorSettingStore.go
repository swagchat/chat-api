package datastore

import (
	"fmt"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
	"github.com/pkg/errors"
)

func rdbCreateOperatorSettingStore(db string) {
	master := RdbStore(db).master()
	tableMap := master.AddTableWithName(model.OperatorSetting{}, tableNameOperatorSetting)
	tableMap.SetKeys(true, "id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating operator setting table. %v.", err))
		return
	}
}

func rdbInsertOperatorSetting(db string, setting *model.OperatorSetting) (*model.OperatorSetting, error) {
	master := RdbStore(db).master()

	if err := master.Insert(setting); err != nil {
		err = errors.Wrap(err, "An error occurred while inserting operator setting")
		logger.Error(err.Error())
		return nil, err
	}

	return setting, nil
}

func rdbSelectOperatorSetting(db, settingID string) (*model.OperatorSetting, error) {
	slave := RdbStore(db).replica()

	var ss []*model.OperatorSetting
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created DESC;", tableNameOperatorSetting)
	_, err := slave.Select(&ss, query)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting operator setting")
		logger.Error(err.Error())
		return nil, err
	}

	if len(ss) == 1 {
		return ss[0], nil
	}

	return nil, nil
}

func rdbUpdateOperatorSetting(db string, setting *model.OperatorSetting) error {
	master := RdbStore(db).master()

	if _, err := master.Update(setting); err != nil {
		err = errors.Wrap(err, "An error occurred while updating operator setting")
		logger.Error(err.Error())
		return err
	}

	return nil
}
