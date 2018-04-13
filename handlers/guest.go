package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setGuestMux() {
	mux.PostFunc("/guests", commonHandler(postGuest))
	mux.GetFunc("/guests/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(getGuest)))
}

func postGuest(w http.ResponseWriter, r *http.Request) {
	var post models.User
	if err := decodeBody(r, &post); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, pd := services.PostGuest(r.Context(), &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getGuest(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	user, pd := services.GetGuest(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}
