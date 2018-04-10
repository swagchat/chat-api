package handlers

import (
	"net/http"

	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetSettingMux() {
	Mux.GetFunc("/setting", colsHandler(datastoreHandler(GetSetting)))
}

func GetSetting(w http.ResponseWriter, r *http.Request) {
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	setting, pd := services.GetSetting(dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", setting)
}
