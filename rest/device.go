package rest

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
)

func setDeviceMux() {
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(getDevices)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(getDevice)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(putDevice)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(deleteDevice)))
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	devices, pd := service.GetDevices(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", devices)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")
	i, err := strconv.ParseInt(bone.GetValue(r, "platform"), 10, 32)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
			Error: err,
		})
	}
	platform := int32(i)

	device, pd := service.GetDevice(r.Context(), userID, platform)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", device)
}

func putDevice(w http.ResponseWriter, r *http.Request) {
	var put model.Device
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.UserID = bone.GetValue(r, "userId")
	i, err := strconv.ParseInt(bone.GetValue(r, "platform"), 10, 32)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
			Error: err,
		})
	}
	put.Platform = int32(i)

	device, errRes := service.PutDevice(r.Context(), &put)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	if device == nil {
		respond(w, r, http.StatusNotModified, "", nil)
	} else {
		respond(w, r, http.StatusOK, "application/json", device)
	}
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")
	i, err := strconv.ParseInt(bone.GetValue(r, "platform"), 10, 32)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
			Error: err,
		})
	}
	platform := int32(i)

	errRes := service.DeleteDevice(r.Context(), userID, platform)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
