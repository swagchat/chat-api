package service_test

import (
	"testing"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	"github.com/fairway-corp/operator-api/utils"
	_ "github.com/mattn/go-sqlite3"
)

func TestGuest(t *testing.T) {
	t.Run("Create guest", func(t *testing.T) {
		metaData := utils.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("failed %s test %v", "Get guest", err)
		}
		req := &model.CreateGuestRequest{}
		req.UserID = "user-id-0001"
		req.Name = "Name"
		req.PictureURL = "http://example.com/dummy.png"
		req.InformationURL = "http://example.com"
		req.MetaData = metaData
		req.Public = true
		req.CanBlock = true
		req.Lang = "ja"
		_, errRes := service.CreateGuest(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Create guest", errRes)
		}
	})
	t.Run("Get guest", func(t *testing.T) {
		req := &model.GetGuestRequest{}
		req.UserID = "user-id-0001"
		res, errRes := service.GetGuest(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Get guest", errRes)
		}
		if res == nil {
			t.Fatalf("failed %s test", "Get guest")
		}
	})
}
