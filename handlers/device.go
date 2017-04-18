package handlers

import (
	"net/http"
	"strconv"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/services"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func SetDeviceMux() {
	basePath := "/users"
	Mux.PostFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$/devices"), ColsHandler(PostDevice))
	Mux.DeleteFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$/devices/#platform^[1-9]$"), ColsHandler(DeleteDevice))
}

func PostDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice models.Device
	if err := decodeBody(r, &requestDevice); err != nil {
		respondJsonDecodeError(w, r, "Create device item")
		return
	}

	userId := bone.GetValue(r, "userId")
	device, pd := services.CreateDevice(userId, &requestDevice)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", device)
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
