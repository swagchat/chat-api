package rest

import (
	"net/http"
	"time"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/go-zoo/bone"
)

func setGuestMux() {
	mux.PostFunc("/guests", commonHandler(postGuest))
	mux.GetFunc("/guests/:userId", commonHandler(getGuest))
}

func postGuest(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGuestRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	res, pd := service.CreateGuest(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    res.AccessToken,
		Domain:   r.Header.Get(utils.HeaderWorkspace),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, jwtCookie)

	respond(w, r, http.StatusCreated, "application/json", res)
}

func getGuest(w http.ResponseWriter, r *http.Request) {
	var req model.GetGuestRequest

	userID := bone.GetValue(r, "userId")
	req.UserID = userID
	res, pd := service.GetGuest(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	jwtCookie := &http.Cookie{
		Name:  "jwt",
		Value: res.AccessToken,
		// Domain:   r.Header.Get(utils.HeaderWorkspace),
		Domain:   "client.fairway-corp.co.jp",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, jwtCookie)

	respond(w, r, http.StatusOK, "application/json", res)
}
