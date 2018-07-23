package rest

import (
	"net/http"

	"github.com/fairway-corp/operator-api/model"

	"github.com/fairway-corp/operator-api/service"
)

func setGuestSettingMux() {
	mux.PostFunc("/guestSettings", commonHandler(postGuestSetting))
	mux.GetFunc("/guestSettings", commonHandler(getGuestSetting))
	mux.PutFunc("/guestSettings", commonHandler(putGuestSetting))
}

func postGuestSetting(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGuestSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	res, pd := service.CreateGuestSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", res)
}

func getGuestSetting(w http.ResponseWriter, r *http.Request) {
	var req model.GetGuestSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	res, pd := service.GetGuestSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", res)
}

func putGuestSetting(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateGuestSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.UpdateGuestSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
