package idp

import (
	"context"

	"github.com/fairway-corp/operator-api/utils"
)

type provider interface {
	Init() error
	Create() (string, string, error)
	GetToken() (string, error)
}

func Provider(ctx context.Context) provider {
	cfg := utils.Config()
	var p provider

	switch cfg.IdP.Provider {
	case "local":
		p = &localProvider{
			ctx: ctx,
		}
	case "keycloak":
		p = &keycloakProvider{
			ctx:          ctx,
			baseEndpoint: cfg.IdP.Keycloak.BaseEndpoint,
		}
	}

	return p
}
