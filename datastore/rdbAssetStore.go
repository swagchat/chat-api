package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateAssetStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateAssetStore")
	defer span.Finish()

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
		logger.Error(fmt.Sprintf("An error occurred while creating asset. %v.", err))
		return
	}
}

func rdbInsertAsset(ctx context.Context, db string, asset *model.Asset) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertAsset")
	defer span.Finish()

	master := RdbStore(db).master()

	if err := master.Insert(asset); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting assett. %v.", err))
		return err
	}

	return nil
}

func rdbSelectAsset(ctx context.Context, db, assetID string) (*model.Asset, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectAsset")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var assets []*model.Asset
	query := fmt.Sprintf("SELECT * FROM %s WHERE asset_id=:assetId AND deleted = 0;", tableNameAsset)
	params := map[string]interface{}{"assetId": assetID}
	_, err := replica.Select(&assets, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting asset. %v.", err))
		return nil, err
	}

	if len(assets) > 0 {
		return assets[0], nil
	}

	return nil, nil
}
