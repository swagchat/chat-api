package rest

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setDeviceMux() {
	mux.PostFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(postDevice)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/devices", commonHandler(selfResourceAuthzHandler(getDevices)))
	mux.PutFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(putDevice)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$/devices/#platform^[1-9]$", commonHandler(selfResourceAuthzHandler(deleteDevice)))
}

func postDevice(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.postDevice")
	defer span.Finish()

	var req model.CreateDeviceRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	errRes := service.CreateDevice(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.getDevices")
	defer span.Finish()

	req := &model.GetDevicesRequest{}
	req.UserID = bone.GetValue(r, "userId")

	devices, errRes := service.GetDevices(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", devices)
}

func putDevice(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.putDevice")
	defer span.Finish()

	var req model.UpdateDeviceRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

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

	req.Platform = int32(i)

	errRes := service.UpdateDevice(r.Context(), &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "application/json", nil)
	// if device == nil {
	// 	respond(w, r, http.StatusNotModified, "", nil)
	// } else {
	// 	respond(w, r, http.StatusOK, "application/json", device)
	// }
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "rest.deleteDevices")
	defer span.Finish()

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

	req.Platform = int32(i)

	errRes := service.DeleteDevice(r.Context(), req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}
