package handlers

import (
	"net/http"
	"net/url"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/go-zoo/bone"
)

func SetRoomMux() {
	Mux.PostFunc("/rooms", ColsHandler(PostRoom))
	Mux.GetFunc("/rooms", ColsHandler(GetRooms))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$", ColsHandler(GetRoom))
	Mux.PutFunc("/rooms/#roomId^[a-z0-9-]$", ColsHandler(PutRoom))
	Mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$", ColsHandler(DeleteRoom))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/messages", ColsHandler(GetRoomMessages))
}

func PostRoom(w http.ResponseWriter, r *http.Request) {
	var post models.Room
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create room item")
		return
	}

	room, pd := services.CreateRoom(&post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", room)
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)
	rooms, pd := services.GetRooms(requestParams)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", rooms)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	room, pd := services.GetRoom(roomId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func PutRoom(w http.ResponseWriter, r *http.Request) {
	var put models.Room
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Update room item")
		return
	}

	roomId := bone.GetValue(r, "roomId")
	room, pd := services.PutRoom(roomId, &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "", room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	pd := services.DeleteRoom(roomId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)
	messages, pd := services.GetRoomMessages(roomId, requestParams)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
