package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/go-zoo/bone"
)

func SetRoomUserMux() {
	Mux.PostFunc("/rooms/#roomId^[a-z0-9-]$/users", ColsHandler(PostRoomUsers))
	Mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", ColsHandler(PutRoomUser))
	Mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", ColsHandler(DeleteRoomUsers))
	Mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users", ColsHandler(PutRoomUsers))
}

func PostRoomUsers(w http.ResponseWriter, r *http.Request) {
	var post models.RequestRoomUserIds
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	roomUsers, pd := services.PostRoomUsers(roomId, &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}

func PutRoomUser(w http.ResponseWriter, r *http.Request) {
	var put models.RoomUser
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update room's user item")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	userId := bone.GetValue(r, "userId")
	roomUser, pd := services.PutRoomUser(roomId, userId, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUser)
}

func DeleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var deleteRus models.RequestRoomUserIds
	if err := decodeBody(r, &deleteRus); err != nil {
		respondJsonDecodeError(w, r, "Deleting room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	roomUsers, pd := services.DeleteRoomUsers(roomId, &deleteRus)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}

func PutRoomUsers(w http.ResponseWriter, r *http.Request) {
	var put models.RequestRoomUserIds
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Adding room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	roomUsers, pd := services.PutRoomUsers(roomId, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}
