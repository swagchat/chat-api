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
	Mux.PutFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(PutRoomUsers))
	Mux.DeleteFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users"), ColsHandler(DeleteRoomUsers))
	Mux.PutFunc(utils.AppendStrings(basePath, "/#roomId^[a-z0-9-]$/users/#userId^[a-z0-9-]$"), ColsHandler(PutRoomUser))
}

func PostRoomUsers(w http.ResponseWriter, r *http.Request) {
	var post models.RoomUsers
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	resRoomUser, pd := services.PostRoomUsers(roomId, &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	status := http.StatusCreated
	if len(resRoomUser.Errors) > 0 {
		status = http.StatusInternalServerError
	}

	respond(w, r, status, "application/json", resRoomUser)
}

func PutRoomUsers(w http.ResponseWriter, r *http.Request) {
	var put models.RoomUsers
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Adding room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	rus, pd := services.PutRoomUsers(roomId, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusInternalServerError, "application/json", rus)
}

func DeleteRoomUsers(w http.ResponseWriter, r *http.Request) {
	var deleteRus models.RoomUsers
	if err := decodeBody(r, &deleteRus); err != nil {
		respondJsonDecodeError(w, r, "Deleting room's user list")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	resRoomUser, pd := services.DeleteRoomUsers(roomId, &deleteRus)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	if len(resRoomUser.Errors) > 0 {
		respond(w, r, http.StatusInternalServerError, "application/json", resRoomUser)
	} else {
		respond(w, r, http.StatusNoContent, "", nil)
	}
}

func PutRoomUser(w http.ResponseWriter, r *http.Request) {
	var put models.RoomUser
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update room's user item")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	userId := bone.GetValue(r, "userId")
	pd := services.PutRoomUser(roomId, userId, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
