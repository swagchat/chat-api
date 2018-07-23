package rest

import (
	"net/http"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	scpb "github.com/swagchat/protobuf"
)

func setWebhookMux() {
	mux.PostFunc("/webhooks/room", commonHandler(postRoomCreationEvent))
	mux.PostFunc("/webhooks/message", commonHandler(postMessageSendEvent))
}

func postRoomCreationEvent(w http.ResponseWriter, r *http.Request) {
	var req model.Room
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.RecvWebhookRoom(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func postMessageSendEvent(w http.ResponseWriter, r *http.Request) {
	var req scpb.Message
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.RecvWebhookMessage(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
