package idp

import (
	"context"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type provider interface {
	Init() error
	Post(ctx context.Context) (*models.User, error)
	Get(ctx context.Context, userID string) (*models.User, error)
}

func Provider() provider {
	cfg := utils.Config()
	var p provider

	switch cfg.IdP.Provider {
	case "local":
		p = &localProvider{}
	case "keycloak":
		p = &keycloakProvider{
			baseEndpoint: cfg.IdP.Keycloak.BaseEndpoint,
		}
	}

	return p
}
