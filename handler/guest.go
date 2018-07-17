package handler

import (
	"net/http"
	"time"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
)

func setGuestMux() {
	mux.PostFunc("/guests", commonHandler(postGuest))
	mux.GetFunc("/guests/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(getGuest)))
}

func postGuest(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGuestRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, pd := service.CreateGuest(r.Context(), &req)
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
	req := &model.GetGuestRequest{}
	req.UserID = bone.GetValue(r, "userId")

	user, pd := service.GetGuest(r.Context(), req)
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
