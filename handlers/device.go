package handlers

import (
	"net/http"
	"strconv"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/go-zoo/bone"
)

func SetDeviceMux() {
	Mux.PostFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", ColsHandler(PostDevice))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", ColsHandler(GetDevices))
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", ColsHandler(GetDevice))
	Mux.PutFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", ColsHandler(PutDevice))
	Mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", ColsHandler(DeleteDevice))
}

func PostDevice(w http.ResponseWriter, r *http.Request) {
	var post models.Device
	if err := decodeBody(r, &post); err != nil {
		respondJsonDecodeError(w, r, "Create device item")
		return
	}

	userId := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	device, pd := services.CreateDevice(userId, platform, &post)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", device)
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, pd := services.GetDevices()
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

	userId := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	device, pd := services.PutDevice(userId, platform, &put)
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
