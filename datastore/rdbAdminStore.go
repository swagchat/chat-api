package datastore

import (
	"log"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
	"time"
)

func RdbCreateAdminStore() {
	tableMap := dbMap.AddTableWithName(models.Admin{}, TABLE_NAME_ADMIN)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
		return
	}
	dRes := RdbSelectLatestAdmin()
	if dRes.Data == nil {
		dRes = RdbInsertAdmin()
		if dRes.ProblemDetail != nil {
			// error
		}
	}
}

func RdbInsertAdmin() StoreResult {
	result := StoreResult{}
	admin := &models.Admin{
		Token:   utils.GenerateToken(utils.TOKEN_LENGTH),
		Created: time.Now().Unix(),
	}
	if err := dbMap.Insert(admin); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating admin item.", err)
	}
	result.Data = admin
	return result
}

func RdbSelectLatestAdmin() StoreResult {
	result := StoreResult{}
	var admins []*models.Admin
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ADMIN, " ORDER BY created DESC LIMIT 1;")
	if _, err := dbMap.Select(&admins, query); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting admin item.", err)
	}
	if len(admins) > 0 {
		result.Data = admins[0]
	}
	return result
}
