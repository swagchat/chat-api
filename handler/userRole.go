package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf"
)

func setUserRoleMux() {
	mux.PostFunc("/userRoles", commonHandler(adminAuthzHandler(postUserRole)))
	mux.GetFunc("/userRoles", commonHandler(adminAuthzHandler(getUserRole)))
	mux.DeleteFunc("/userRoles", commonHandler(adminAuthzHandler(deleteUserRole)))
}

func postUserRole(w http.ResponseWriter, r *http.Request) {
	var req scpb.CreateUserRolesRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.CreateUserRoles(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getUserRole(w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	userID := ""
	if userIDSli, ok := params["userId"]; ok {
		userID = userIDSli[0]
	}

	roleID := ""
	var roleIDint32 int32
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

	if userID != "" {
		req := &scpb.GetRoleIdsOfUserRoleRequest{
			UserId: userID,
		}
		roleIDs, pd := service.GetRoleIDsOfUserRole(r.Context(), req)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		respond(w, r, http.StatusOK, "application/json", roleIDs)
		return
	}

	if roleIDint32 > 0 {
		req := &scpb.GetUserIdsOfUserRoleRequest{
			RoleId: roleIDint32,
		}
		userIDs, pd := service.GetUserIDsOfUserRole(r.Context(), req)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		respond(w, r, http.StatusOK, "application/json", userIDs)
		return
	}

	respond(w, r, http.StatusNotFound, "application/json", nil)
}

func deleteUserRole(w http.ResponseWriter, r *http.Request) {
	var req scpb.DeleteUserRoleRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.DeleteUserRole(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", nil)
}
