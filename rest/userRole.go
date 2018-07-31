package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setUserRoleMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(selfResourceAuthzHandler(postUserRole)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(selfResourceAuthzHandler(deleteUserRole)))
}

func postUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postUserRole")
	defer span.Finish()

	var req model.CreateUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.CreateUserRoles(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func deleteUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteUserRole")
	defer span.Finish()

	var req model.DeleteUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.DeleteUserRoles(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
