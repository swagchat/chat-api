package handlers

import (
	"net/http"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
	"github.com/go-zoo/bone"
)

func SetRoomUserMux() {
	Mux.PutFunc(utils.AppendStrings("/", utils.API_VERSION, "/rooms/#roomId^[a-z0-9-]$/users"), colsHandler(PutRoomUsers))
	Mux.PutFunc(utils.AppendStrings("/", utils.API_VERSION, "/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$"), colsHandler(PutRoomUser))
	Mux.DeleteFunc(utils.AppendStrings("/", utils.API_VERSION, "/rooms/#roomId^[a-z0-9-]$/users"), colsHandler(DeleteRoomUsers))
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

func PutRoomUser(w http.ResponseWriter, r *http.Request) {
	var put models.RoomUser
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update room's user item")
		return
	}

	put.RoomId = bone.GetValue(r, "roomId")
	put.UserId = bone.GetValue(r, "userId")
	roomUser, pd := services.PutRoomUser(&put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, roomUser.Modified)
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
