package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetUserMux() {
	Mux.PostFunc("/users", colsHandler(jwtHandler(datastoreHandler(datastoreHandler(PostUser)))))
	Mux.GetFunc("/users", colsHandler(datastoreHandler(datastoreHandler(GetUsers))))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(datastoreHandler(GetUser))))
	Mux.GetFunc("/profiles/#userId^[a-z0-9-]$", colsHandler(contactsAuthHandler(datastoreHandler(GetProfile))))
	Mux.PutFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(datastoreHandler(PutUser))))
	Mux.DeleteFunc("/users/#userId^[a-z0-9-]$", colsHandler(userAuthHandler(datastoreHandler(DeleteUser))))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/unreadCount", colsHandler(userAuthHandler(datastoreHandler(GetUserUnreadCount))))
}

func userAuthHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := bone.GetValue(r, "userId")
		sub := r.Header.Get(jwtSub)
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
		sub := r.Header.Get(jwtSub)
		if userID != "" && sub != "" {
			dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)
			pd := services.ContactsAuth(userID, sub, dsCfg)
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

	jwt := r.Context().Value(ctxJwt).(*models.JWT)
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	user, pd := services.PostUser(&post, jwt, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	users, pd := services.GetUsers(dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	user, pd := services.GetUser(userId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	user, pd := services.GetProfile(userId, dsCfg)
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
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	user, pd := services.PutUser(&put, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	pd := services.DeleteUser(userId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func GetUserUnreadCount(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	userUnreadCount, pd := services.GetUserUnreadCount(userId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", userUnreadCount)
}
