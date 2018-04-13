package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setRoomMux() {
	mux.PostFunc("/rooms", commonHandler(postRoom))
	mux.GetFunc("/rooms", commonHandler(adminAuthzHandler(getRooms)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(getRoom)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoom)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(deleteRoom)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/messages", commonHandler(roomMemberAuthzHandler(getRoomMessages)))
}

func postRoom(w http.ResponseWriter, r *http.Request) {
	var post models.Room
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	room, pd := services.PostRoom(r.Context(), &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", room)
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)

	rooms, pd := services.GetRooms(r.Context(), requestParams)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", rooms)
}

func getRoom(w http.ResponseWriter, r *http.Request) {
	roomID := bone.GetValue(r, "roomId")

	room, pd := services.GetRoom(r.Context(), roomID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func putRoom(w http.ResponseWriter, r *http.Request) {
	var put models.Room
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.RoomId = bone.GetValue(r, "roomId")

	room, pd := services.PutRoom(r.Context(), &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func deleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := bone.GetValue(r, "roomId")

	pd := services.DeleteRoom(r.Context(), roomID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func getRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomID := bone.GetValue(r, "roomId")
	params, _ := url.ParseQuery(r.URL.RawQuery)

	messages, pd := services.GetRoomMessages(r.Context(), roomID, params)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
