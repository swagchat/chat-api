package handler

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessages)))
	mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessages(w http.ResponseWriter, r *http.Request) {
	var post model.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	mRes := service.PostMessage(r.Context(), &post)
	// if len(mRes.MessageIds) == 0 {
	// 	respond(w, r, mRes.Errors[0].Status, "application/json", mRes)
	// 	return
	// }
	if len(mRes.Errors) > 0 {
		pd := mRes.Errors[0]
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", mRes)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	messageID := bone.GetValue(r, "messageId")

	message, pd := service.GetMessage(r.Context(), messageID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, message.Modified)
	respond(w, r, http.StatusOK, "application/json", message)
}
