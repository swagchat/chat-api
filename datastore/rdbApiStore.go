package datastore

import (
	"log"
	"strconv"

	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateApiStore() {
	tableMap := dbMap.AddTableWithName(models.Api{}, TABLE_NAME_API)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
		return
	}
	dRes := RdbSelectLatestApi("admin")
	if dRes.Data == nil {
		dRes = RdbInsertApi("admin")
		if dRes.ProblemDetail != nil {
			// error
		}
	}
}

func RdbInsertApi(name string) StoreResult {
	result := StoreResult{}
	api := &models.Api{
		Name:    name,
		Key:     utils.CreateApiKey(),
		Secret:  utils.GenerateToken(utils.TOKEN_LENGTH),
		Created: time.Now().Unix(),
		Expired: 0,
	}
	if err := dbMap.Insert(api); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating api item.", err)
	}
	result.Data = api
	return result
}

func RdbSelectLatestApi(name string) StoreResult {
	result := StoreResult{}
	var apis []*models.Api
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_API, " WHERE name=:name AND (expired=0 OR expired>", nowTimestampString, ") ORDER BY created DESC LIMIT 1;")
	log.Println(query)
	params := map[string]interface{}{"name": name}
	if _, err := dbMap.Select(&apis, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting api item.", err)
	}
	if len(apis) > 0 {
		result.Data = apis[0]
	}
	return result
}
