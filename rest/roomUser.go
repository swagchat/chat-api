package rest

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
	"github.com/swagchat/chat-api/utils"
)

func setRoomUserMux() {
	mux.PostFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(postRoomUsers)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(getRoomUsers)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoomUser)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$/users", commonHandler(roomMemberAuthzHandler(deleteRoomUsers)))
}

func postRoomUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postRoomUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.AddRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	errRes := service.AddRoomUsers(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", nil)
}

func getRoomUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getRoomUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	params, _ := url.ParseQuery(r.URL.RawQuery)

	req := &model.RetrieveRoomUsersRequest{}
	req.RoomID = bone.GetValue(r, "roomId")

	commaSeparatedRoleIDs := ""
	if commaSeparatedRoleIDsSli, ok := params["roleIds"]; ok {
		commaSeparatedRoleIDs = commaSeparatedRoleIDsSli[0]
	}

	var roleIDs []int32
	if commaSeparatedRoleIDs != "" {
		roleIDs = utils.CommaSeparatedStringToInt32(commaSeparatedRoleIDs)
		if len(roleIDs) > 0 {
			req.RoleIDs = roleIDs
		}
	}

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		roomUserIDs, errRes := service.RetrieveRoomUserIDs(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", roomUserIDs)
	} else {
		roomUsers, errRes := service.RetrieveRoomUsers(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", roomUsers)
	}
}

func putRoomUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("putRoomUser", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.UpdateRoomUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")
	req.UserID = bone.GetValue(r, "userId")

	errRes := service.UpdateRoomUser(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("deleteRoomUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.DeleteRoomUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	errRes := service.DeleteRoomUsers(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
