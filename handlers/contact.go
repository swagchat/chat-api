package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/services"
	"github.com/go-zoo/bone"
)

func SetContactMux() {
	Mux.GetFunc("/contacts/#userId^[a-z0-9-]$", colsHandler(aclHandler(GetContacts)))
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	contacts, pd := services.GetContacts(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}
