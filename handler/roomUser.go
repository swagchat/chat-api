package handler

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf"
)

func setRoomUserMux() {
	mux.PostFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(postRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func postRoomUsers(w http.ResponseWriter, r *http.Request) {
	var req scpb.CreateRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	pd := service.CreateRoomUsers(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	var req scpb.UpdateRoomUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")
	req.UserID = bone.GetValue(r, "userId")

	pd := service.UpdateRoomUser(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var req scpb.DeleteRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	pd := service.DeleteRoomUsers(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
