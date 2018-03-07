package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func SetUserMux() {
	Mux.PostFunc("/users", colsHandler(aclHandler(PostUser)))
	Mux.GetFunc("/users", colsHandler(GetUsers))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(GetUser)))
	Mux.GetFunc("/profiles/#userId^[a-z0-9-]$", colsHandler(contactsAuthHandler(GetProfile)))
	Mux.PutFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(PutUser)))
	Mux.DeleteFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(DeleteUser)))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/unreadCount", colsHandler(userAuthHandler(GetUserUnreadCount)))
}

func userAuthHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := bone.GetValue(r, "userId")
		sub := r.Header.Get("X-Sub")
		if userID != "" && sub != "" {
			pd := services.UserAuth(userID, sub)
			if pd != nil {
				respondErr(w, r, pd.Status, pd)
				return
			}
		}
		fn(w, r)
	}
}

func contactsAuthHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := bone.GetValue(r, "userId")
		sub := r.Header.Get("X-Sub")
		if userID != "" && sub != "" {
			pd := services.ContactsAuth(userID, sub)
			if pd != nil {
				respondErr(w, r, pd.Status, pd)
				return
			}
		}
		fn(w, r)
	}
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var post models.User
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create user item")
		return
	}

	jwt := &models.JWT{
		Sub: r.Header.Get("X-Sub"),
	}

	user, pd := services.PostUser(&post, jwt)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, pd := services.GetUsers()
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	user, pd := services.GetUser(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	user, pd := services.GetProfile(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	var put models.User
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update user item")
		return
	}

	put.UserId = bone.GetValue(r, "userId")

	user, pd := services.PutUser(&put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	pd := services.DeleteUser(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func GetUserUnreadCount(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	userUnreadCount, pd := services.GetUserUnreadCount(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", userUnreadCount)
}
