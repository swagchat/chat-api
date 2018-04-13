package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setDeviceMux() {
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(getDevices)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(getDevice)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(putDevice)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(deleteDevice)))
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	devices, pd := services.GetDevices(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", devices)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))

	device, pd := services.GetDevice(r.Context(), userID, platform)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", device)
}

func putDevice(w http.ResponseWriter, r *http.Request) {
	var put models.Device
	if err := decodeBody(r, &put); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	put.UserId = bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	put.Platform = platform

	device, pd := services.PutDevice(r.Context(), &put)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
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
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))

	pd := services.DeleteDevice(r.Context(), userID, platform)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
