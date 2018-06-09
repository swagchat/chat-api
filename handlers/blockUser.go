package handlers

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
)

func setBlockUserMux() {
	mux.GetFunc("/users/#userId^[a-z0-9-]$/blocks", commonHandler(selfResourceAuthzHandler(getBlockUsers)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/blocks", commonHandler(selfResourceAuthzHandler(putBlockUsers)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/blocks", commonHandler(selfResourceAuthzHandler(deleteBlockUsers)))
}

func getBlockUsers(w http.ResponseWriter, r *http.Request) {
	userID := bone.GetValue(r, "userId")

	blockUsers, pd := services.GetBlockUsers(r.Context(), userID)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}

func putBlockUsers(w http.ResponseWriter, r *http.Request) {
	var reqUIDs models.RequestBlockUserIDs
	if err := decodeBody(r, &reqUIDs); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")

	blockUsers, pd := services.PutBlockUsers(r.Context(), userID, &reqUIDs)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}

func deleteBlockUsers(w http.ResponseWriter, r *http.Request) {
	var reqUIDs models.RequestBlockUserIDs
	if err := decodeBody(r, &reqUIDs); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	userID := bone.GetValue(r, "userId")

	blockUsers, pd := services.DeleteBlockUsers(r.Context(), userID, &reqUIDs)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusOK, "application/json", blockUsers)
}
