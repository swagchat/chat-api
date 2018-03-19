package datastore

import (
	"log"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateAssetStore() {
	master := RdbStoreInstance().master()
	tableMap := master.AddTableWithName(models.Asset{}, TABLE_NAME_ASSET)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
		return
	}
}

func RdbInsertAsset(asset *models.Asset) StoreResult {
	master := RdbStoreInstance().master()
	result := StoreResult{}
	if err := master.Insert(asset); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating asset item.", err)
	}
	result.Data = asset
	return result
}

func RdbSelectAsset(assetId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var assets []*models.Asset
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ASSET, " WHERE asset_id=:assetId AND deleted = 0;")
	params := map[string]interface{}{"assetId": assetId}
	if _, err := slave.Select(&assets, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting asset item.", err)
	}
	if len(assets) > 0 {
		result.Data = assets[0]
	}
	return result
}
