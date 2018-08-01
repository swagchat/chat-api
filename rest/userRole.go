package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setUserRoleMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(selfResourceAuthzHandler(postUserRole)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(selfResourceAuthzHandler(putUserRole)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(selfResourceAuthzHandler(deleteUserRole)))
}

func postUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postUserRole", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.CreateUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.CreateUserRoles(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func putUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("putUserRole", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.AddUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.AddUserRoles(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("deleteUserRole", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.DeleteUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.DeleteUserRoles(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
