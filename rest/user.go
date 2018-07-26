package rest

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
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
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postUser")
	defer span.Finish()

	var req model.CreateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, errRes := service.CreateUser(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getUsers")
	defer span.Finish()

	req := &model.GetUsersRequest{}
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, nil)
		return
	}

	limit, offset, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	users, errRes := service.GetUsers(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getUser")
	defer span.Finish()

	req := &model.GetUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, errRes := service.GetUser(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.putUser")
	defer span.Finish()

	var req model.UpdateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	user, errRes := service.UpdateUser(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteUser")
	defer span.Finish()

	req := &model.DeleteUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.DeleteUser(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getContacts")
	defer span.Finish()

	req := &model.GetContactsRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, nil)
		return
	}

	limit, offset, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	contacts, errRes := service.GetContacts(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getProfile")
	defer span.Finish()

	req := &model.GetProfileRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, errRes := service.GetProfile(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
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
