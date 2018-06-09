package handlers

import (
	"fmt"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func configValidate() {
	cfg := utils.Config()
	if cfg.IdP.Provider == "keycloak" {
		if cfg.IdP.Keycloak.BaseEndpoint == "" {
			logging.Log(zapcore.FatalLevel, &logging.AppLog{
				Kind:  "storage",
				Error: fmt.Errorf("keycloak base endpoint is empty"),
			})
		}
	}
}
