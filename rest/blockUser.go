package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setBlockUserMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(postBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(getBlockUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blockedUsers", commonHandler(selfResourceAuthzHandler(getBlockedUsers)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(putBlockUsers)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/blockUsers", commonHandler(selfResourceAuthzHandler(deleteBlockUsers)))
}

func postBlockUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postBlockUsers")
	defer span.Finish()

	var req model.CreateBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.CreateBlockUsers(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getBlockUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getBlockUsers")
	defer span.Finish()

	req := &model.GetBlockUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		req.ResponseType = scpb.ResponseType_UserIdList
	} else {
		req.ResponseType = scpb.ResponseType_UserList
	}

	blockUsers, errRes := service.GetBlockUsers(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}

func getBlockedUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getBlockedUsers")
	defer span.Finish()

	req := &model.GetBlockedUsersRequest{}
	req.UserID = bone.GetValue(r, "userId")

	responseType := bone.GetValue(r, "responseType")
	if responseType == "UserIdList" {
		req.ResponseType = scpb.ResponseType_UserIdList
	} else {
		req.ResponseType = scpb.ResponseType_UserList
	}

	blockedUsers, errRes := service.GetBlockedUsers(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockedUsers)
}

func putBlockUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.putBlockUsers")
	defer span.Finish()

	var req model.AddBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.AddBlockUsers(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteBlockUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteBlockUsers")
	defer span.Finish()

	var req model.DeleteBlockUsersRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	errRes := service.DeleteBlockUsers(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
