package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func SetRoomMux() {
	Mux.PostFunc("/rooms", colsHandler(PostRoom))
	Mux.GetFunc("/rooms", colsHandler(GetRooms))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(GetRoom)))
	Mux.PutFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(PutRoom)))
	Mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(DeleteRoom)))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/messages", colsHandler(roomAuthHandler(GetRoomMessages)))
}

func roomAuthHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := bone.GetValue(r, "roomId")
		sub := r.Header.Get("X-Sub")
		if roomID != "" && sub != "" {
			pd := services.RoomAuth(roomID, sub)
			if pd != nil {
				respondErr(w, r, pd.Status, pd)
				return
			}
		}
		fn(w, r)
	}
}

func PostRoom(w http.ResponseWriter, r *http.Request) {
	var post models.Room
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create room item")
		return
	}

	room, pd := services.PostRoom(&post)
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

	put.RoomId = bone.GetValue(r, "roomId")

	room, pd := services.PutRoom(&put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
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
	params, _ := url.ParseQuery(r.URL.RawQuery)
	messages, pd := services.GetRoomMessages(roomId, params)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
