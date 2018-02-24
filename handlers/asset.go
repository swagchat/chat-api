package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/utils"
)

func SetAssetMux() {
	Mux.PostFunc("/assets", colsHandler(PostAsset))
	Mux.GetFunc("/assets/#assetId^[a-z0-9-]$", GetAsset)
}

func PostAsset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		problemDetail := &models.ProblemDetail{
			Title:     "MultipartForm parse error. (Create asset item)",
			ErrorName: models.ERROR_NAME_INVALID_JSON,
		}
		respondErr(w, r, http.StatusBadRequest, problemDetail)
		return
	}

	file, header, err := r.FormFile("asset")
	if err != nil {
		problemDetail := &models.ProblemDetail{
			Title:     "Form value error. (Create asset item)",
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "asset",
					Reason: "asset is required, but it's empty.",
				},
			},
		}
		respondErr(w, r, http.StatusBadRequest, problemDetail)
		return
	}
	defer file.Close()

	var extension string
	contentType := header.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	default:
		problemDetail := &models.ProblemDetail{
			Title:     "Unsupported file. (Create asset item)",
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "Content-Type",
					Reason: utils.AppendStrings("Content-Type is image/jpeg or image/png, but it's ", contentType),
				},
			},
		}
		respondErr(w, r, http.StatusBadRequest, problemDetail)
		return
	}

	storageProvider := storage.GetProvider()
	assetId := utils.CreateUuid()
	assetInfo := &storage.AssetInfo{
		FileName: utils.AppendStrings(assetId, extension),
		Data:     file,
	}
	sourceUrl, problemDetail := storageProvider.Post(assetInfo)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	asset := &models.Asset{
		AssetId:   assetId,
		SourceUrl: sourceUrl,
		Mime:      header.Header.Get("Content-Type"),
	}

	respond(w, r, http.StatusCreated, "application/json", asset)
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
	assetId := bone.GetValue(r, "assetId")

	storageProvider := storage.GetProvider()
	assetInfo := &storage.AssetInfo{
		FileName: assetId,
	}

	bytes, problemDetail := storageProvider.Get(assetInfo)
	if problemDetail != nil {
		respondErr(w, r, problemDetail.Status, problemDetail)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Set("Content-Type", http.DetectContentType(bytes))
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
