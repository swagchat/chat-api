package rest

import (
	"net/http"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setMessageMux() {
	mux.PostFunc("/messages", commonHandler(updateLastAccessedHandler(postMessages)))
	mux.GetFunc("/messages/#messageId^[a-z0-9-]$", commonHandler(updateLastAccessedHandler(getMessage)))
}

func postMessages(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postMessages")
	defer span.Finish()

	var post model.Messages
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	mRes := service.CreateMessages(r.Context(), &post)
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
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getMessage")
	defer span.Finish()

	messageID := bone.GetValue(r, "messageId")

	message, pd := service.GetMessage(r.Context(), messageID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, message.Modified)
	respond(w, r, http.StatusOK, "application/json", message)
}
