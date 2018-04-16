package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setRoomUserMux() {
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(putRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func putRoomUsers(w http.ResponseWriter, r *http.Request) {
	var put models.RequestRoomUserIDs
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	roomID := bone.GetValue(r, "roomId")

	roomUsers, pd := services.PutRoomUsers(r.Context(), roomID, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	var put models.RoomUser
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.RoomID = bone.GetValue(r, "roomId")
	put.UserID = bone.GetValue(r, "userId")

	roomUser, pd := services.PutRoomUser(r.Context(), &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, roomUser.Modified)
	respond(w, r, http.StatusOK, "application/json", roomUser)
}

func deleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var deleteRus models.RequestRoomUserIDs
	if err := decodeBody(r, &deleteRus); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	roomID := bone.GetValue(r, "roomId")

	roomUsers, pd := services.DeleteRoomUsers(r.Context(), roomID, &deleteRus)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}
