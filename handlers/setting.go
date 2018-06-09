package handlers

import (
	"net/http"

	"github.com/swagchat/chat-api/services"
)

func setSettingMux() {
	mux.GetFunc("/setting", commonHandler(getSetting))
}

func getSetting(w http.ResponseWriter, r *http.Request) {
	setting, pd := services.GetSetting(r.Context())
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", setting)
}
