package handlers

import (
	"net/http"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
	"github.com/go-zoo/bone"
)

func SetBlockUserMux() {
	Mux.GetFunc(utils.AppendStrings("/", utils.API_VERSION, "/users/#userId^[a-z0-9-]$/blocks"), colsHandler(GetBlockUsers))
	Mux.PutFunc(utils.AppendStrings("/", utils.API_VERSION, "/users/#userId^[a-z0-9-]$/blocks"), colsHandler(PutBlockUsers))
	Mux.DeleteFunc(utils.AppendStrings("/", utils.API_VERSION, "/users/#userId^[a-z0-9-]$/blocks"), colsHandler(DeleteBlockUsers))
}

func GetBlockUsers(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	blockUsers, pd := services.GetBlockUsers(userId)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}

func PutBlockUsers(w http.ResponseWriter, r *http.Request) {
	var reqUIDs models.RequestBlockUserIds
	if err := decodeBody(r, &reqUIDs); err != nil {
		respondJsonDecodeError(w, r, "Adding block user list")
		return
	}

	userId := bone.GetValue(r, "userId")
	blockUsers, pd := services.PutBlockUsers(userId, &reqUIDs)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}

func DeleteBlockUsers(w http.ResponseWriter, r *http.Request) {
	var reqUIDs models.RequestBlockUserIds
	if err := decodeBody(r, &reqUIDs); err != nil {
		respondJsonDecodeError(w, r, "Deleting block user list")
		return
	}

	userId := bone.GetValue(r, "userId")
	blockUsers, pd := services.DeleteBlockUsers(userId, &reqUIDs)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}
