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
	Mux.PostFunc(utils.AppendStrings(basePath, "/#userId^[a-z0-9-]$/devices/#platform^[1-9]$"), ColsHandler(PostDevice))
}

func PostDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice models.Device
	if err := decodeBody(r, &requestDevice); err != nil {
		respondJsonDecodeError(w, r, "Create user item")
		return
	}

	userId := bone.GetValue(r, "userId")
	platform, _ := strconv.Atoi(bone.GetValue(r, "platform"))
	device, pd := services.CreateDevice(userId, platform, &requestDevice)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", device)
}
