package rest

import (
	"net/http"
	"net/url"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setUserRoleMux() {
	mux.PostFunc("/userRoles", commonHandler(adminAuthzHandler(postUserRole)))
	mux.GetFunc("/userRoles", commonHandler(adminAuthzHandler(getUserRole)))
	mux.DeleteFunc("/userRoles", commonHandler(adminAuthzHandler(deleteUserRole)))
}

func postUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postUserRole")
	defer span.Finish()

	var req model.CreateUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	errRes := service.CreateUserRoles(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getUserRole")
	defer span.Finish()

	params, _ := url.ParseQuery(r.URL.RawQuery)

	userID := ""
	if userIDSli, ok := params["userId"]; ok {
		userID = userIDSli[0]
	}

	roleID := ""
	if roleIDSli, ok := params["roleId"]; ok {
		roleID = roleIDSli[0]
	}

	if userID != "" && roleID != "" {
		respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*model.InvalidParam{
				&model.InvalidParam{
					Name:   "userId, roleId",
					Reason: "Be sure to specify either userId or roleId.",
				},
			},
			Status: http.StatusBadRequest,
		})
		return
	}

	if userID != "" {
		gRIDsReq := &model.GetRoleIdsOfUserRoleRequest{}
		gRIDsReq.UserID = userID
		roleIDs, errRes := service.GetRoleIDsOfUserRole(r.Context(), gRIDsReq)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}

		respond(w, r, http.StatusOK, "application/json", roleIDs)
		return
	}

	var roleIDint32 int32
	if roleID != "" {
		i, err := strconv.ParseInt(roleID, 10, 32)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "roleId",
						Reason: "roleId must be numeric.",
					},
				},
				Status: http.StatusBadRequest,
			})
			return
		}
		roleIDint32 = int32(i)
	}

	if roleIDint32 > 0 {
		gUIDsReq := &model.GetUserIdsOfUserRoleRequest{}
		gUIDsReq.RoleID = roleIDint32
		userIDs, errRes := service.GetUserIDsOfUserRole(r.Context(), gUIDsReq)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}

		respond(w, r, http.StatusOK, "application/json", userIDs)
		return
	}

	respond(w, r, http.StatusNotFound, "application/json", nil)
}

func deleteUserRole(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteUserRole")
	defer span.Finish()

	var req model.DeleteUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	errRes := service.DeleteUserRoles(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNotFound, "application/json", nil)
}
