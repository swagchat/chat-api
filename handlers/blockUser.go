package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func SetBlockUserMux() {
	Mux.GetFunc("/users/#userId^[a-z0-9-]$/blocks", colsHandler(userAuthHandler(datastoreHandler(GetBlockUsers))))
	Mux.PutFunc("/users/#userId^[a-z0-9-]$/blocks", colsHandler(userAuthHandler(datastoreHandler(PutBlockUsers))))
	Mux.DeleteFunc("/users/#userId^[a-z0-9-]$/blocks", colsHandler(userAuthHandler(datastoreHandler(DeleteBlockUsers))))
}

func GetBlockUsers(w http.ResponseWriter, r *http.Request) {
	userId := bone.GetValue(r, "userId")
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	blockUsers, pd := services.GetBlockUsers(userId, dsCfg)
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
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	blockUsers, pd := services.PutBlockUsers(userId, &reqUIDs, dsCfg)
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
	dsCfg := r.Context().Value(ctxDsCfg).(*utils.Datastore)

	blockUsers, pd := services.DeleteBlockUsers(userId, &reqUIDs, dsCfg)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}
