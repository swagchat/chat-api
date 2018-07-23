package service_test

import (
	"testing"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	_ "github.com/mattn/go-sqlite3"
)

func TestGuestSetting(t *testing.T) {
	t.Run("Get guest setting", func(t *testing.T) {
		req := &model.GetGuestSettingRequest{}
		res, errRes := service.GetGuestSetting(ctx, req)
		if res != nil {
			t.Fatalf("failed %s test", "Get guest setting")
		}
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Get guest setting", errRes)
		}
	})
	t.Run("Create guest setting", func(t *testing.T) {
		falseValue := false
		req := &model.CreateGuestSettingRequest{}
		req.EnableWebchat = &falseValue
		res, errRes := service.CreateGuestSetting(ctx, req)
		if res == nil {
			t.Fatalf("failed %s test %v", "Create guest setting", errRes)
		}
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Create guest setting", errRes)
		}
	})
	t.Run("Get guest setting", func(t *testing.T) {
		req := &model.GetGuestSettingRequest{}
		res, errRes := service.GetGuestSetting(ctx, req)
		if res == nil {
			t.Fatalf("failed %s test", "Get guest setting")
		}
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Get guest setting", errRes)
		}
	})
	t.Run("Update guest setting", func(t *testing.T) {
		trueValue := true
		req := &model.UpdateGuestSettingRequest{}
		req.EnableWebchat = &trueValue
		errRes := service.UpdateGuestSetting(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Update guest setting", errRes)
		}
	})
	t.Run("Get guest setting", func(t *testing.T) {
		req := &model.GetGuestSettingRequest{}
		res, errRes := service.GetGuestSetting(ctx, req)
		if res == nil {
			t.Fatalf("failed %s test %v", "Get guest setting", errRes)
		}
		if errRes != nil {
			t.Fatalf("failed %s test %v", "Get guest setting", errRes)
		}
	})
}
