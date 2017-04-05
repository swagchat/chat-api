package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func SetMessageMux() {
	basePath := "/messages"

	Mux.PostFunc(basePath, ColsHandler(PostMessage))
	Mux.GetFunc(utils.AppendStrings(basePath, "/#messageId^[a-z0-9-]$"), ColsHandler(GetMessage))
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	var requestMessages models.Messages
	if err := decodeBody(r, &requestMessages); err != nil {
		respondJsonDecodeError(w, r, "Create message item")
		return
	}

	messages, problemDetail := services.CreateMessage(&requestMessages)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", messages)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	messageId := bone.GetValue(r, "messageId")
	message, problemDetail := services.GetMessage(messageId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", message)
}
