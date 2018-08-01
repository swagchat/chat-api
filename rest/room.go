package rest

import (
	"net/http"
	"net/url"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setRoomMux() {
	mux.PostFunc("/rooms", commonHandler(postRoom))
	mux.GetFunc("/rooms", commonHandler(adminAuthzHandler(getRooms)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(getRoom)))
	mux.PutFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(putRoom)))
	mux.DeleteFunc("/rooms/#roomId^[a-z0-9-]$", commonHandler(roomMemberAuthzHandler(deleteRoom)))
	mux.GetFunc("/rooms/#roomId^[a-z0-9-]$/messages", commonHandler(roomMemberAuthzHandler(updateLastAccessedHandler(getRoomMessages))))
}

func postRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postRoom", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.CreateRoomRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	room, errRes := service.CreateRoom(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", room)
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getRooms", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.GetRoomsRequest{}
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	limit, offset, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	rooms, errRes := service.GetRooms(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", rooms)
}

func getRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getRoom", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.GetRoomRequest{}

	roomID := bone.GetValue(r, "roomId")
	req.RoomID = roomID

	room, errRes := service.GetRoom(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func putRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("putRoom", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var req model.UpdateRoomRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.RoomID = bone.GetValue(r, "roomId")

	room, errRes := service.UpdateRoom(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", room)
}

func deleteRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("deleteRoom", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.DeleteRoomRequest{}

	roomID := bone.GetValue(r, "roomId")
	req.RoomID = roomID

	errRes := service.DeleteRoom(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func getRoomMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getRoomMessages", "rest")
	defer tracer.Provider(ctx).Finish(span)

	req := &model.GetRoomMessagesRequest{}

	roomID := bone.GetValue(r, "roomId")
	req.RoomID = roomID

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	limit, offset, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	messages, errRes := service.GetRoomMessages(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", messages)
}
