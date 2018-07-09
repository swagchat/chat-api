package handler

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/service"
)

func setRoomUserMux() {
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(putRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func putRoomUsers(w http.ResponseWriter, r *http.Request) {
	var req protobuf.PostRoomUserReq
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	roomUsers, pd := service.PutRoomUsers(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	var put protobuf.RoomUser
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.RoomID = bone.GetValue(r, "roomId")
	put.UserID = bone.GetValue(r, "userId")

	roomUser, pd := service.PutRoomUser(r.Context(), &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUser)
}

func deleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var req protobuf.DeleteRoomUserReq
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	roomUsers, pd := service.DeleteRoomUsers(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}
