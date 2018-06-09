package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func rdbCreateAssetStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Asset{}, tableNameAsset)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create asset table error",
			Error:   err,
		})
	}
}

func rdbInsertAsset(db string, asset *models.Asset) (*models.Asset, error) {
	master := RdbStore(db).master()

	if err := master.Insert(asset); err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating asset")
	}

	return asset, nil
}

func rdbSelectAsset(db, assetID string) (*models.Asset, error) {
	replica := RdbStore(db).replica()

	var assets []*models.Asset
	query := utils.AppendStrings("SELECT * FROM ", tableNameAsset, " WHERE asset_id=:assetId AND deleted = 0;")
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
