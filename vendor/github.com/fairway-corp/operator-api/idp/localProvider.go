package idp

import (
	"context"

	"github.com/fairway-corp/operator-api/utils"
)

type localProvider struct {
	ctx context.Context
}

func (lp *localProvider) Init() error {
	return nil
}

func (lp *localProvider) Create() (string, string, error) {
	// pd := req.Validate()
	// if pd != nil {
	// 	return nil, pd.Error
	// }

	// // Create user
	// gimei := gimei.NewName()
	// user := req.GenerateUser()
	// user.Name = fmt.Sprintf("%s(%s)(仮)", gimei.Kanji(), gimei.Katakana())
	// req.UserID = user.UserID
	// userRoles := req.GenerateUserRoles()
	// user, err := datastore.Provider(ctx).InsertUser(user, datastore.UserOptionInsertRoles(userRoles))
	// if err != nil {
	// 	return nil, errors.Wrap(err, "")
	// }

	return utils.GenerateUUID(), "", nil
}

func (lp *localProvider) GetToken() (string, error) {
	return "", nil
}
