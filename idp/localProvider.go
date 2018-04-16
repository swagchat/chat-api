package idp

import (
	"context"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type localProvider struct {
}

func (lp *localProvider) Init() error {
	return nil
}

func (lp *localProvider) Post(ctx context.Context) (*models.User, error) {
	uuid := utils.GenerateUUID()
	user := &models.User{
		UserID: uuid,
		Name:   uuid,
	}

	user.BeforeInsertGuest()

	user, err := datastore.Provider(ctx).InsertUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return user, nil
}

func (lp *localProvider) Get(ctx context.Context, userID string) (*models.User, error) {
	user, err := datastore.Provider(ctx).SelectUser(userID, true, true, true)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return user, nil
}
