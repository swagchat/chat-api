package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func SetDeviceMux() {
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", colsHandler(GetDevices))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", colsHandler(GetDevice))
	Mux.PutFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", colsHandler(PutDevice))
	Mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", colsHandler(DeleteDevice))
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	devices, pd := services.GetDevices(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", devices)
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	device, pd := services.GetDevice(userId, platform)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", device)
}

func PutDevice(w http.ResponseWriter, r *http.Request) {
	var put models.Device
	if err := decodeBody(r, &put); err != nil {
		respondJsonDecodeError(w, r, "Create device item")
		return
	}

	put.UserId = bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	put.Platform = platform
	device, pd := services.PutDevice(&put)
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

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	pd := services.DeleteDevice(userId, platform)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
