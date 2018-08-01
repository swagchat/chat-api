package rest

import (
	"net/http"

	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setSettingMux() {
	mux.GetFunc("/setting", commonHandler(getSetting))
}

func getSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getSetting", "rest")
	defer tracer.Provider(ctx).Finish(span)

	setting, errRes := service.GetSetting(ctx)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", setting)
}
