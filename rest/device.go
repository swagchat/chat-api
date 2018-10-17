package rest

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/betchi/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setDeviceMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(postDevice)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(getDevices)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(deleteDevice)))
}

func postDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "postDevice", "rest")
	defer tracer.Finish(span)

	var req model.AddDeviceRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	device, errRes := service.AddDevice(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	if device == nil {
		respond(w, r, http.StatusNotModified, "", nil)
	} else {
		respond(w, r, http.StatusCreated, "application/json", device)
	}
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getDevices", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveDevicesRequest{}
	req.UserID = bone.GetValue(r, "userId")

	devices, errRes := service.RetrieveDevices(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", devices)
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "deleteDevice", "rest")
	defer tracer.Finish(span)

	req := &model.DeleteDeviceRequest{}
	req.UserID = bone.GetValue(r, "userId")

	i, err := strconv.ParseInt(bone.GetValue(r, "platform"), 10, 32)
	if err != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "platform",
				Reason: "Platform must be numeric.",
			},
		}
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		respondError(w, r, errRes)
		return
	}

	req.Platform = scpb.Platform(i)

	errRes := service.DeleteDevice(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
