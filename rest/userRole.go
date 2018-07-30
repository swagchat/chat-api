package rest

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setUserRoleMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(adminAuthzHandler(postUserRole)))
	mux.GetFunc("/roles/#roleId^[0-9]$/userIds", commonHandler(adminAuthzHandler(getUserIDsOfUserRole)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/roles", commonHandler(adminAuthzHandler(deleteUserRole)))
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

func getUserIDsOfUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getUserIDsOfUserRole")
	defer span.Finish()

	req := &model.GetUserIdsOfUserRoleRequest{}

	roleIDString := bone.GetValue(r, "roleId")
	roleIDInt, err := strconv.ParseInt(roleIDString, 10, 32)
	if err != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roleId",
				Reason: "roleId must be numeric.",
			},
		}
		errRes := model.NewErrorResponse("Failed to get userIds of user role.", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		respondError(w, r, errRes)
		return
	}

	req.RoleID = int32(roleIDInt)
	userIDs, errRes := service.GetUserIDsOfUserRole(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}
	respond(w, r, http.StatusOK, "application/json", userIDs)
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
