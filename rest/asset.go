package rest

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/betchi/tracer"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setAssetMux() {
	mux.PostFunc("/assets", commonHandler(postAsset))
	mux.GetFunc("/assets/:filename", commonHandler(getAsset))
	mux.GetFunc("/assets/:filename/info", commonHandler(getAssetInfo))
}

func postAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "postAsset", "rest")
	defer tracer.Finish(span)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		errRes := model.NewErrorResponse("MultipartForm parse error.", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	file, header, err := r.FormFile("asset")
	if err != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "asset",
				Reason: "asset is required, but it's empty.",
			},
		}
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		respondError(w, r, errRes)
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

	asset, errRes := service.PostAsset(ctx, contentType, file, size, width, height)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", asset)
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getAsset", "rest")
	defer tracer.Finish(span)

	filename := bone.GetValue(r, "filename")
	assetID := utils.GetFileNameWithoutExt(filename)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	bytes, asset, errRes := service.GetAsset(ctx, assetID, ifModifiedSince)
	if errRes != nil {
		respondError(w, r, errRes)
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
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getAssetInfo", "rest")
	defer tracer.Finish(span)

	filename := bone.GetValue(r, "filename")
	assetID := utils.GetFileNameWithoutExt(filename)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	asset, errRes := service.GetAssetInfo(ctx, assetID, ifModifiedSince)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", asset)
}
