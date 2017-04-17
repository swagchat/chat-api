package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func SetUserMux() {
	basePath := "/users"
	Mux.PostFunc(basePath, ColsHandler(PostUser))
	Mux.GetFunc(basePath, ColsHandler(GetUsers))
	Mux.GetFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$"), ColsHandler(GetUser))
	Mux.PutFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$"), ColsHandler(PutUser))
	Mux.DeleteFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$"), ColsHandler(DeleteUser))
	//	Mux.GetFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$/rooms"), ColsHandler(GetUserRooms))
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var post models.User
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create user item")
		return
	}

	user, pd := services.CreateUser(&post)
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

func PutUser(w http.ResponseWriter, r *http.Request) {
	var put models.User
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update user item")
		return
	}

	userId := bone.GetValue(r, "userId")
	user, pd := services.PutUser(userId, &put)
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

/*
func GetUserRooms(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	user, pd := services.GetUserRooms(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "", user)
}
*/
