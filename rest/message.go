package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessages)))
	mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("postMessages", "rest")
	defer tracer.Provider(ctx).Finish(span)

	var post model.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	mRes := service.CreateMessages(ctx, &post)
	// if len(mRes.MessageIds) == 0 {
	// 	respond(w, r, mRes.Errors[0].Status, "application/json", mRes)
	// 	return
	// }
	if len(mRes.Errors) > 0 {
		errRes := mRes.Errors[0]
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", mRes)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.Provider(ctx).StartSpan("getMessage", "rest")
	defer tracer.Provider(ctx).Finish(span)

	messageID := bone.GetValue(r, "messageId")

	message, errRes := service.GetMessage(ctx, messageID)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	setLastModified(w, message.Modified)
	respond(w, r, http.StatusOK, "application/json", message)
}
