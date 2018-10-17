package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/betchi/tracer"
)

func setBlockUserMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(postBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(getBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockedUsers", commonHandler(selfResourceAuthzHandler(getBlockedUsers)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(deleteBlockUsers)))
}

func postBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "postBlockUsers", "rest")
	defer tracer.Finish(span)

	var req model.AddBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.AddBlockUsers(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", nil)
}

func getBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getBlockUsers", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveBlockUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		blockUserIDs, errRes := service.RetrieveBlockUserIDs(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockUserIDs)
	} else {
		blockUsers, errRes := service.RetrieveBlockUsers(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockUsers)
	}
}

func getBlockedUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getBlockedUsers", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveBlockedUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		blockedUserIDs, errRes := service.RetrieveBlockedUserIDs(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockedUserIDs)
	} else {
		blockedUsers, errRes := service.RetrieveBlockedUsers(ctx, req)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		respond(w, r, http.StatusOK, "application/json", blockedUsers)
	}
}

func deleteBlockUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "deleteBlockUsers", "rest")
	defer tracer.Finish(span)

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
