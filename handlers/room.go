package handlers

import (
	"net/http"
	"net/url"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func SetRoomMux() {
	basePath := "/rooms"
	Mux.PostFunc(basePath, ColsHandler(PostRoom))
	Mux.GetFunc(basePath, ColsHandler(GetRooms))
	Mux.GetFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$"), ColsHandler(GetRoom))
	Mux.PutFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$"), ColsHandler(PutRoom))
	Mux.DeleteFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$"), ColsHandler(DeleteRoom))

	Mux.GetFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/messages"), ColsHandler(GetRoomMessages))
}

func PostRoom(w http.ResponseWriter, r *http.Request) {
	var requestRoom models.Room
	if err := decodeBody(r, &requestRoom); err != nil {
		respondJsonDecodeError(w, r, "Create room item")
		return
	}

	room, problemDetail := services.CreateRoom(&requestRoom)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", room)
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)
	rooms, problemDetail := services.GetRooms(requestParams)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", rooms)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	room, problemDetail := services.GetRoom(roomId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func PutRoom(w http.ResponseWriter, r *http.Request) {
	var requestRoom models.Room
	if err := decodeBody(r, &requestRoom); err != nil {
		respondJsonDecodeError(w, r, "Update room item")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	room, problemDetail := services.PutRoom(roomId, &requestRoom)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "", room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	resRoomUser, problemDetail := services.DeleteRoom(roomId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	if len(resRoomUser.Errors) > 0 {
		respond(w, r, http.StatusInternalServerError, "application/json", resRoomUser)
	} else {
		respond(w, r, http.StatusNoContent, "", nil)
	}
}

func GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)
	messages, problemDetail := services.GetRoomMessages(roomId, requestParams)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
