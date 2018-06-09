package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func setAssetMux() {
	mux.PostFunc("/assets", commonHandler(postAsset))
	mux.GetFunc("/assets/:filename", commonHandler(getAsset))
	mux.GetFunc("/assets/:filename/info", commonHandler(getAssetInfo))
}

func postAsset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:     "MultipartForm parse error. (Create asset item)",
			ErrorName: models.ERROR_NAME_INVALID_JSON,
		}
		respondErr(w, r, http.StatusBadRequest, pd)
		return
	}

	file, header, err := r.FormFile("asset")
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "asset",
					Reason: "asset is required, but it's empty.",
				},
			},
		}
		respondErr(w, r, http.StatusBadRequest, pd)
		return
	}
	defer file.Close()

	contentType := r.FormValue("mime")
	if contentType == "" {
		contentType = header.Header.Get("Content-Type")
	}
	size := header.Size
	width, _ := strconv.Atoi(r.FormValue("width"))
	height, _ := strconv.Atoi(r.FormValue("height"))

	asset, pd := services.PostAsset(r.Context(), contentType, file, size, width, height)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", asset)
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	filename := bone.GetValue(r, "filename")
	assetID := utils.GetFileNameWithoutExt(filename)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	bytes, asset, pd := services.GetAsset(r.Context(), assetID, ifModifiedSince)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	if asset == nil {
		respondErr(w, r, http.StatusNotFound, nil)
		return
	}

	setLastModified(w, asset.Modified)
	// w.Header().Set("Cache-Control", "max-age:86400, public")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Set("Content-Type", http.DetectContentType(bytes))
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func getAssetInfo(w http.ResponseWriter, r *http.Request) {
	filename := bone.GetValue(r, "filename")
	assetID := utils.GetFileNameWithoutExt(filename)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	asset, pd := services.GetAssetInfo(r.Context(), assetID, ifModifiedSince)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	if asset == nil {
		respondErr(w, r, http.StatusNotFound, nil)
		return
	}

	respond(w, r, http.StatusOK, "application/json", asset)
}
