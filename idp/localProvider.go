package idp

import (
	"context"
	"fmt"

	"github.com/betchi/go-gimei"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

type localProvider struct {
}

func (lp *localProvider) Init() error {
	return nil
}

func (lp *localProvider) Post(ctx context.Context) (*model.User, error) {
	uuid := utils.GenerateUUID()
	gimei := gimei.NewName()
	user := &model.User{
		UserID: uuid,
		Name:   fmt.Sprintf("%s(%s)(仮)", gimei.Kanji(), gimei.Katakana()),
	}

	user.BeforeInsertGuest()

	general := &model.UserRole{
		UserID: user.UserID,
		RoleID: utils.RoleGeneral,
	}
	guest := &model.UserRole{
		UserID: user.UserID,
		RoleID: utils.RoleGuest,
	}
	roles := []*model.UserRole{general, guest}
	user, err := datastore.Provider(ctx).InsertUser(user, roles)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return user, nil
}

func (lp *localProvider) Get(ctx context.Context, userID string) (*model.User, error) {
	user, err := datastore.Provider(ctx).SelectUser(userID, datastore.WithBlocks(true), datastore.WithDevices(true), datastore.WithRooms(true))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return user, nil
}
