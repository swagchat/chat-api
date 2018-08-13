package rest

import (
	"net/http"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessage)))
	// mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postMessage", "rest")
	defer tracer.Provider(ctx).Finish(span)

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
// 	span := tracer.Provider(ctx).StartSpan("getMessage", "rest")
// 	defer tracer.Provider(ctx).Finish(span)

// 	messageID := bone.GetValue(r, "messageId")

// 	message, errRes := service.GetMessage(ctx, messageID)
// 	if errRes != nil {
// 		respondError(w, r, errRes)
// 		return
// 	}

// 	setLastModified(w, message.Modified)
// 	respond(w, r, http.StatusOK, "application/json", message)
// }
