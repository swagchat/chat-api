package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetContactMux() {
	Mux.GetFunc("/contacts/#userId^[a-z0-9-]$", colsHandler(datastoreHandler(GetContacts)))
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	contacts, pd := services.GetContacts(userId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}
