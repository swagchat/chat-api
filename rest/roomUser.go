package rest

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setRoomUserMux() {
	mux.PostFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(postRoomUsers)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(getRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func postRoomUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postRoomUsers")
	defer span.Finish()

	var req model.CreateRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	errRes := service.CreateRoomUsers(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getRoomUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getRoomUsers")
	defer span.Finish()

	params, _ := url.ParseQuery(r.URL.RawQuery)

	req := &model.GetRoomUsersRequest{}
	req.RoomID = bone.GetValue(r, "roomId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		req.ResponseType = scpb.ResponseType_UserIdList
	} else {
		req.ResponseType = scpb.ResponseType_UserList
	}

	commaSeparatedRoleIDs := ""
	if commaSeparatedRoleIDsSli, ok := params["roleIds"]; ok {
		commaSeparatedRoleIDs = commaSeparatedRoleIDsSli[0]
	}

	var roleIDs []int32
	if commaSeparatedRoleIDs != "" {
		roleIDs = utils.CommaSeparatedStringsToInt32(commaSeparatedRoleIDs)
		if len(roleIDs) > 0 {
			req.RoleIDs = roleIDs
		}
	}

	roomUsers, errRes := service.GetRoomUsers(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}
	respond(w, r, http.StatusOK, "application/json", roomUsers)
	return
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.putRoomUser")
	defer span.Finish()

	var req model.UpdateRoomUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")
	req.UserID = bone.GetValue(r, "userId")

	errRes := service.UpdateRoomUser(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteRoomUsers")
	defer span.Finish()

	var req model.DeleteRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	errRes := service.DeleteRoomUsers(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
