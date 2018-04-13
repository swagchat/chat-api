package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setUserMux() {
	mux.PostFunc("/users", commonHandler(postUser))
	mux.GetFunc("/users", commonHandler(adminAuthzHandler(getUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(getUser)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(putUser)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(deleteUser)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/unreadCount", commonHandler(selfResourceAuthzHandler(getUserUnreadCount)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/contacts", commonHandler(selfResourceAuthzHandler(getContacts)))
	mux.GetFunc("/profiles/#userId^[a-z0-9-]$", commonHandler(contactsAuthzHandler(getProfile)))
}

func postUser(w http.ResponseWriter, r *http.Request) {
	var post models.User
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, pd := services.PostUser(r.Context(), &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, pd := services.GetUsers(r.Context())
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	user, pd := services.GetUser(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	var put models.User
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.UserId = bone.GetValue(r, "userId")

	user, pd := services.PutUser(r.Context(), &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	pd := services.DeleteUser(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func getUserUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	userUnreadCount, pd := services.GetUserUnreadCount(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", userUnreadCount)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	contacts, pd := services.GetContacts(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	user, pd := services.GetProfile(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}
