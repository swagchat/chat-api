package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetRoomMux() {
	Mux.PostFunc("/rooms", colsHandler(datastoreHandler(PostRoom)))
	Mux.GetFunc("/rooms", colsHandler(datastoreHandler(GetRooms)))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(datastoreHandler(GetRoom))))
	Mux.PutFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(datastoreHandler(PutRoom))))
	Mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$", colsHandler(roomAuthHandler(datastoreHandler(DeleteRoom))))
	Mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/messages", colsHandler(roomAuthHandler(datastoreHandler(GetRoomMessages))))
}

func roomAuthHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := bone.GetValue(r, "roomId")
		sub := r.Header.Get(jwtSub)
		dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

		if roomID != "" && sub != "" {
			pd := services.RoomAuth(roomID, sub, dsCfg)
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

	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	room, pd := services.PostRoom(&post, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", room)
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	requestParams, _ := url.ParseQuery(r.URL.RawQuery)
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	rooms, pd := services.GetRooms(requestParams, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", rooms)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	room, pd := services.GetRoom(roomId, dsCfg)
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
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	room, pd := services.PutRoom(&put, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	pd := services.DeleteRoom(roomId, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomId := bone.GetValue(r, "roomId")
	params, _ := url.ParseQuery(r.URL.RawQuery)
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	messages, pd := services.GetRoomMessages(roomId, params, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
