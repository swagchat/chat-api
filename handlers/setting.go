package handlers

import (
	"net/http"

	"github.com/swagchat/chat-api/services"
)

func SetSettingMux() {
	Mux.GetFunc("/setting", colsHandler(GetSetting))
}

func GetSetting(w http.ResponseWriter, r *http.Request) {
	setting, pd := services.GetSetting()
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", setting)
}
