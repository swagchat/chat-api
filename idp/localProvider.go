package idp

import (
	"context"
	"fmt"

	"github.com/betchi/go-gimei"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

type localProvider struct {
}

func (lp *localProvider) Init() error {
	return nil
}

func (lp *localProvider) Post(ctx context.Context, req *model.CreateGuestRequest) (*model.User, error) {
	pd := req.Validate()
	if pd != nil {
		return nil, pd.Error
	}

	// Create user
	gimei := gimei.NewName()
	user := req.GenerateUser()
	user.Name = fmt.Sprintf("%s(%s)(仮)", gimei.Kanji(), gimei.Katakana())
	req.UserID = user.UserID
	userRoles := req.GenerateUserRoles()
	user, err := datastore.Provider(ctx).InsertUser(user, datastore.UserOptionInsertRoles(userRoles))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return user, nil
}

func (lp *localProvider) Get(ctx context.Context, req *model.GetGuestRequest) (*model.User, error) {
	user, err := datastore.Provider(ctx).SelectUser(req.UserID, datastore.UserOptionWithBlocks(true), datastore.UserOptionWithDevices(true), datastore.UserOptionWithRooms(true))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return user, nil
}
