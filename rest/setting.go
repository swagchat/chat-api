package rest

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/service"
)

func setSettingMux() {
	mux.GetFunc("/setting", commonHandler(getSetting))
}

func getSetting(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getSetting")
	defer span.Finish()

	setting, errRes := service.GetSetting(r.Context())
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", setting)
}
