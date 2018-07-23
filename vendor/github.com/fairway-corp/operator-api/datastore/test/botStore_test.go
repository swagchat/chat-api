package datastore_test

import (
	"testing"
	"time"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/utils"
	_ "github.com/mattn/go-sqlite3"
)

func TestBotStore(t *testing.T) {
	t.Run("Insert bot", func(t *testing.T) {
		falseValue := false
		nowTimestamp := time.Now().Unix()
		bot := &model.Bot{}
		bot.BotID = utils.GenerateUUID()
		bot.UserID = "userId"
		bot.Service = "service"
		bot.ProjectID = "projectId"
		bot.ServiceAccount = "{}"
		bot.Suggest = &falseValue
		bot.Created = nowTimestamp
		bot.Modified = nowTimestamp
		bot.Deleted = int64(0)

		res, err := datastore.Provider(ctx).InsertBot(bot)
		if res == nil {
			t.Fatalf("failed %s test %v", "Insert guest setting", err)
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Insert guest setting", err)
		}
	})
	t.Run("Select bot", func(t *testing.T) {
		res, err := datastore.Provider(ctx).SelectBot("userId")
		if res == nil {
			t.Fatalf("failed %s test", "Select bot")
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Select bot", err)
		}
	})
}
