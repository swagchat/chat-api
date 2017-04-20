package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/go-zoo/bone"
)

func SetMessageMux() {
	Mux.PostFunc("/messages", ColsHandler(PostMessage))
	Mux.GetFunc("/messages/#messageId^[a-z0-9-]$", ColsHandler(GetMessage))
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	var post models.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create message item")
		return
	}

	messages, pd := services.CreateMessage(&post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", messages)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	messageId := bone.GetValue(r, "messageId")
	message, pd := services.GetMessage(messageId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", message)
}
