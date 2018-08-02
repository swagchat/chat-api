package datastore

import (
	"context"
	"fmt"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
)

func rdbCreateAssetStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateAssetStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.Asset{}, tableNameAsset)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating asset table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertAsset(ctx context.Context, dbMap *gorp.DbMap, asset *model.Asset) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertAsset", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	if err := dbMap.Insert(asset); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting assett. %v.", err))
		return err
	}

	return nil
}

func rdbSelectAsset(ctx context.Context, dbMap *gorp.DbMap, assetID string) (*model.Asset, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectAsset", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var assets []*model.Asset
	query := fmt.Sprintf("SELECT * FROM %s WHERE asset_id=:assetId AND deleted = 0;", tableNameAsset)
	params := map[string]interface{}{"assetId": assetID}
	_, err := dbMap.Select(&assets, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting asset. %v.", err))
		return nil, err
	}

	if len(assets) > 0 {
		return assets[0], nil
	}

	return nil, nil
}
