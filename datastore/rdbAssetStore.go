package datastore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateAssetStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Asset{}, tableNameAsset)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertAsset(db string, asset *model.Asset) (*model.Asset, error) {
	master := RdbStore(db).master()

	if err := master.Insert(asset); err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating asset")
	}

	return asset, nil
}

func rdbSelectAsset(db, assetID string) (*model.Asset, error) {
	replica := RdbStore(db).replica()

	var assets []*model.Asset
	query := fmt.Sprintf("SELECT * FROM %s WHERE asset_id=:assetId AND deleted = 0;", tableNameAsset)
	params := map[string]interface{}{"assetId": assetID}
	_, err := replica.Select(&assets, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting asset")
	}

	if len(assets) > 0 {
		return assets[0], nil
	}

	return nil, nil
}
