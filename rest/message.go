package rest

import (
	"net/http"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/betchi/tracer"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessage)))
	// mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "postMessage", "rest")
	defer tracer.Finish(span)

	var req model.SendMessageRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	message, errRes := service.SendMessage(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", message)
}

// func getMessage(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	span := tracer.StartSpan(ctx, "getMessage", "rest")
// 	defer tracer.Finish(span)

// 	messageID := bone.GetValue(r, "messageId")

// 	message, errRes := service.GetMessage(ctx, messageID)
// 	if errRes != nil {
// 		respondError(w, r, errRes)
// 		return
// 	}

// 	setLastModified(w, message.Modified)
// 	respond(w, r, http.StatusOK, "application/json", message)
// }
