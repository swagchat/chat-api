package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessages)))
	mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessages(w http.ResponseWriter, r *http.Request) {
	var post models.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	mRes := services.PostMessage(r.Context(), &post)
	if len(mRes.MessageIds) == 0 {
		respond(w, r, mRes.Errors[0].Status, "application/json", mRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", mRes)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	messageID := bone.GetValue(r, "messageId")

	message, pd := services.GetMessage(r.Context(), messageID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, message.Modified)
	respond(w, r, http.StatusOK, "application/json", message)
}
