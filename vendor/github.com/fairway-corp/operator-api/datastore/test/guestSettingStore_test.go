package datastore_test

import (
	"testing"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/model"
	_ "github.com/mattn/go-sqlite3"
)

func TestGuestSettingStore(t *testing.T) {
	t.Run("Insert guest setting", func(t *testing.T) {
		falseValue := false
		gs := &model.GuestSetting{}
		gs.EnableWebchat = &falseValue
		res, err := datastore.Provider(ctx).InsertGuestSetting(gs)
		if res == nil {
			t.Fatalf("failed %s test %v", "Insert guest setting", err)
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Insert guest setting", err)
		}
	})
	t.Run("Select guest setting", func(t *testing.T) {
		res, err := datastore.Provider(ctx).SelectGuestSetting()
		if res == nil {
			t.Fatalf("failed %s test", "Select guest setting")
		}
		if err != nil {
			t.Fatalf("failed %s test %v", "Select guest setting", err)
		}
	})
	t.Run("Update guest setting", func(t *testing.T) {
		trueValue := true
		gs := &model.GuestSetting{}
		gs.EnableWebchat = &trueValue
		err := datastore.Provider(ctx).UpdateGuestSetting(gs)
		if err != nil {
			t.Fatalf("failed %s test %v", "Update guest setting", err)
		}
	})
}
