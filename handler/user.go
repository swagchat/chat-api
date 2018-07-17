package handler

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setUserMux() {
	mux.PostFunc("/users", commonHandler(adminAuthzHandler(postUser)))
	mux.GetFunc("/users", commonHandler(adminAuthzHandler(getUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(getUser)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(putUser)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(deleteUser)))

	// mux.GetFunc("/users/#userId^[a-z0-9-]$/unreadCount", commonHandler(selfResourceAuthzHandler(getUserUnreadCount)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/contacts", commonHandler(selfResourceAuthzHandler(getContacts)))
	mux.GetFunc("/profiles/#userId^[a-z0-9-]$", commonHandler(contactsAuthzHandler(getProfile)))
}

func postUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, pd := service.CreateUser(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	req := &model.GetUsersRequest{}
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, nil)
		return
	}

	limit, offset, order, pd := setPagingParams(params)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}
	req.Limit = limit
	req.Offset = offset
	req.Order = order

	users, pd := service.GetUsers(r.Context(), req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	req := &model.GetUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, pd := service.GetUser(r.Context(), req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	user, pd := service.UpdateUser(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	req := &model.DeleteUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	pd := service.DeleteUser(r.Context(), req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

// func getUserUnreadCount(w http.ResponseWriter, r *http.Request) {
// 	userID := bone.GetValue(r, "userId")

// 	userUnreadCount, pd := service.GetUserUnreadCount(r.Context(), userID)
// 	if pd != nil {
// 		respondErr(w, r, pd.Status, pd)
// 		return
// 	}

// 	respond(w, r, http.StatusOK, "application/json", userUnreadCount)
// }

func getContacts(w http.ResponseWriter, r *http.Request) {
	req := &model.GetContactsRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, nil)
		return
	}

	limit, offset, order, pd := setPagingParams(params)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}
	req.Limit = limit
	req.Offset = offset
	req.Order = order

	contacts, pd := service.GetContacts(r.Context(), req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	req := &model.GetProfileRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, pd := service.GetProfile(r.Context(), req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}
