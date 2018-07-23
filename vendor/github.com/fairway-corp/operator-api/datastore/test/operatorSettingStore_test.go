package datastore_test

import (
	"testing"
	"time"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/model"
	_ "github.com/mattn/go-sqlite3"
)

func TestOperatorSettingStore(t *testing.T) {
	t.Run("Insert operator setting", func(t *testing.T) {
		nowTimestamp := time.Now().Unix()
		os := &model.OperatorSetting{}
		os.SettingID = "SettingID"
		os.SiteID = "SiteID"
		os.Domain = "example.com"
		os.SystemUserID = "SystemUserID"
		// os.FirstMessage =
		// os.TimeoutMessage =
		os.OperatorBaseURL = "http://example.com"
		os.NotificationSlackURL = "http://example.slack.com"
		os.Created = nowTimestamp
		os.Modified = nowTimestamp
		res, err := datastore.Provider(ctx).InsertOperatorSetting(os)
		if res == nil {
			t.Fatalf("failed %s test %v", "Insert operator setting", err)
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Insert operator setting", err)
		}
	})
	t.Run("Select operator setting", func(t *testing.T) {
		res, err := datastore.Provider(ctx).SelectOperatorSetting("SettingID")
		if res == nil {
			t.Fatalf("failed %s test", "Select operator setting")
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Select operator setting", err)
		}
	})
	t.Run("Update operator setting", func(t *testing.T) {
		os := &model.OperatorSetting{}
		os.SiteID = "SiteID-update"
		err := datastore.Provider(ctx).UpdateOperatorSetting(os)
		if err != nil {
			t.Fatalf("failed %s test %v", "Update operator setting", err)
		}
	})
}
