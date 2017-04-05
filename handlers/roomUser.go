package handlers

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func SetRoomUserMux() {
	basePath := "/rooms"
	Mux.PostFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(PostRoomUsers))
	//	Mux.GetFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(GetRoomUsers))
	Mux.PutFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(PutRoomUsers))
	Mux.DeleteFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(DeleteRoomUsers))
	Mux.PutFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$"), ColsHandler(PutRoomUser))
}

func PostRoomUsers(w http.ResponseWriter, r *http.Request) {
	var requestRoomUsers models.RoomUsers
	if err := decodeBody(r, &requestRoomUsers); err != nil {
		respondJsonDecodeError(w, r, "Create room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	resRoomUser, problemDetail := services.PostRoomUsers(roomId, &requestRoomUsers)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	status := http.StatusCreated
	if len(resRoomUser.Errors) > 0 {
		status = http.StatusInternalServerError
	}

	respond(w, r, status, "application/json", resRoomUser)
}

/*
func GetRoomUsers(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	room, problemDetail := services.GetRoomUsers(roomId)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}
*/
func PutRoomUsers(w http.ResponseWriter, r *http.Request) {
	var requestRoomUsers models.RoomUsers
	if err := decodeBody(r, &requestRoomUsers); err != nil {
		respondJsonDecodeError(w, r, "Adding room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	resRoomUser, problemDetail := services.PutRoomUsers(roomId, &requestRoomUsers)
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

func DeleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var requestRoomUsers models.RoomUsers
	if err := decodeBody(r, &requestRoomUsers); err != nil {
		respondJsonDecodeError(w, r, "Deleting room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	resRoomUser, problemDetail := services.DeleteRoomUsers(roomId, &requestRoomUsers)
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

func PutRoomUser(w http.ResponseWriter, r *http.Request) {
	var requestRoomUser models.RoomUser
	if err := decodeBody(r, &requestRoomUser); err != nil {
		respondJsonDecodeError(w, r, "Update room's user item")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	userId := bone.GetValue(r, "userId")
	problemDetail := services.PutRoomUser(roomId, userId, &requestRoomUser)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
