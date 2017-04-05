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
	var requestUser models.User
	if err := decodeBody(r, &requestUser); err != nil {
		respondJsonDecodeError(w, r, "Create user item")
		return
	}

	user, problemDetail := services.CreateUser(&requestUser)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, problemDetail := services.GetUsers()
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	user, problemDetail := services.GetUser(userId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	var requestUser models.User
	if err := decodeBody(r, &requestUser); err != nil {
		respondJsonDecodeError(w, r, "Update user item")
		return
	}

	userId := bone.GetValue(r, "userId")
	user, problemDetail := services.PutUser(userId, &requestUser)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "", user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	resRoomUser, problemDetail := services.DeleteUser(userId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	if len(resRoomUser.Errors) > 0 {
		respond(w, r, http.StatusInternalServerError, "application/json", resRoomUser)
	} else {
		respond(w, r, http.StatusNoContent, "", nil)
	}
}

/*
func GetUserRooms(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	user, problemDetail := services.GetUserRooms(userId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "", user)
}
*/
