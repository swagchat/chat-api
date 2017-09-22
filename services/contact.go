// Business Logic

package services

import (
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

func GetContacts(userId string) (*models.Users, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectContacts(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	contacts := &models.Users{
		Users: dRes.Data.([]*models.User),
	}
	return contacts, nil
}
