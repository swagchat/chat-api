package datastore

import (
	"log"
	"strconv"

	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateSettingStore() {
	master := RdbStoreInstance().master()
	tableMap := master.AddTableWithName(models.Setting{}, TABLE_NAME_SETTING)
	tableMap.SetKeys(true, "id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
		return
	}
}

func RdbSelectLatestSetting() StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var settings []*models.Setting
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SETTING, " WHERE expired=0 OR expired>", nowTimestampString, " ORDER BY created DESC LIMIT 1;")
	if _, err := slave.Select(&settings, query, nil); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting setting item.", err)
	}
	if len(settings) > 0 {
		result.Data = settings[0]
	}
	return result
}
