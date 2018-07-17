package handler

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
)

func setRoomUserMux() {
	mux.PostFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(postRoomUsers)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(getRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func postRoomUsers(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRoomUsersRequest
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

func getRoomUsers(w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	gUIDsReq := &model.GetUserIdsOfRoomUserRequest{}
	gUIDsReq.RoomID = bone.GetValue(r, "roomId")

	commaSeparatedRoleIDs := ""
	if commaSeparatedRoleIDsSli, ok := params["roleIds"]; ok {
		commaSeparatedRoleIDs = commaSeparatedRoleIDsSli[0]
	}

	var roleIDs []int32
	if commaSeparatedRoleIDs != "" {
		roleIDs = utils.CommaSeparatedStringsToInt32(commaSeparatedRoleIDs)
		if len(roleIDs) > 0 {
			gUIDsReq.RoleIDs = roleIDs
		}
	}

	if len(roleIDs) > 0 {
		userIDs, pd := service.GetUserIDsOfRoomUser(r.Context(), gUIDsReq)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		respond(w, r, http.StatusOK, "application/json", userIDs)
		return
	}

	respond(w, r, http.StatusNotFound, "application/json", nil)
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateRoomUserRequest
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
	var req model.DeleteRoomUsersRequest
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
