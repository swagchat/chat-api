package handlers

import (
	"net/http"
	"time"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
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

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Domain:   r.Header.Get(utils.HeaderWorkspace),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, jwtCookie)

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getGuest(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	user, pd := services.GetGuest(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Domain:   r.Header.Get(utils.HeaderWorkspace),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, jwtCookie)

	respond(w, r, http.StatusOK, "application/json", user)
}
