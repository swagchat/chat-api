// Business Logic

package services

import (
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func GetContacts(userId string, dsCfg *utils.Datastore) (*models.Users, *models.ProblemDetail) {
	contacts, err := datastore.Provider(dsCfg).SelectContacts(userId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get contact list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.Users{
		Users: contacts,
	}, nil
}
