package datastore

import (
	"log"
	"strconv"

	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateApiStore() {
	master := RdbStoreInstance().Master()
	tableMap := master.AddTableWithName(models.Api{}, TABLE_NAME_API)
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
	dRes := RdbSelectLatestApi("admin")
	if dRes.Data == nil {
		dRes = RdbInsertApi("admin")
		if dRes.ProblemDetail != nil {
			// error
		}
	}
}

func RdbInsertApi(name string) StoreResult {
	master := RdbStoreInstance().Master()
	result := StoreResult{}
	api := &models.Api{
		Name:    name,
		Key:     utils.CreateApiKey(),
		Secret:  utils.GenerateToken(utils.TOKEN_LENGTH),
		Created: time.Now().Unix(),
		Expired: 0,
	}
	if err := master.Insert(api); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating api item.", err)
	}
	result.Data = api
	return result
}

func RdbSelectLatestApi(name string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var apis []*models.Api
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_API, " WHERE name=:name AND (expired=0 OR expired>", nowTimestampString, ") ORDER BY created DESC LIMIT 1;")
	params := map[string]interface{}{"name": name}
	if _, err := slave.Select(&apis, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting api item.", err)
	}
	if len(apis) > 0 {
		result.Data = apis[0]
	}
	return result
}
