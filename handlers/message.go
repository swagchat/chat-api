package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetMessageMux() {
	Mux.PostFunc("/messages", colsHandler(datastoreHandler(PostMessages)))
	Mux.GetFunc("/messages/#messageId^[a-z0-9-]$", colsHandler(datastoreHandler(GetMessage)))
}

func PostMessages(w http.ResponseWriter, r *http.Request) {
	var post models.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create message item")
		return
	}

	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	mRes := services.PostMessage(&post, dsCfg)
	if len(mRes.MessageIds) == 0 {
		respond(w, r, mRes.Errors[0].Status, "application/json", mRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", mRes)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	messageId := bone.GetValue(r, "messageId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	message, pd := services.GetMessage(messageId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, message.Modified)
	respond(w, r, http.StatusOK, "application/json", message)
}
