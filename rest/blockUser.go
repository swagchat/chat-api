package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setBlockUserMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(postBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(getBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockedUsers", commonHandler(selfResourceAuthzHandler(getBlockedUsers)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(deleteBlockUsers)))
}

func postBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postBlockUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.CreateBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.CreateBlockUsers(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getBlockUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.GetBlockUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		blockUserIDs, errRes := service.GetBlockUserIDs(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockUserIDs)
	} else {
		blockUsers, errRes := service.GetBlockUsers(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockUsers)
	}
}

func getBlockedUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getBlockedUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.GetBlockedUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		blockedUserIDs, errRes := service.GetBlockedUserIDs(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockedUserIDs)
	} else {
		blockedUsers, errRes := service.GetBlockedUsers(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockedUsers)
	}
}

func deleteBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("deleteBlockUsers", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.DeleteBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.DeleteBlockUsers(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
