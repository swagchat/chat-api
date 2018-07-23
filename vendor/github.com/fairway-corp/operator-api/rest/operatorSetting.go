package rest

import (
	"net/http"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
)

func setOperatorSettingMux() {
	mux.PostFunc("/guestSettings", commonHandler(postOperatorSetting))
	mux.GetFunc("/guestSettings", commonHandler(getOperatorSetting))
	mux.PutFunc("/guestSettings", commonHandler(putOperatorSetting))
}

func postOperatorSetting(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOperatorSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	res, pd := service.CreateOperatorSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", res)
}

func getOperatorSetting(w http.ResponseWriter, r *http.Request) {
	var req model.GetOperatorSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	res, pd := service.GetOperatorSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", res)
}

func putOperatorSetting(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateOperatorSettingRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	pd := service.UpdateOperatorSetting(r.Context(), &req)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}
