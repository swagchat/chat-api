package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
)

var (
	mux            *bone.Mux
	allowedMethods = []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"PATCH",
		"DELETE",
	}
	noBodyStatusCodes = []int{
		http.StatusNotFound,
		http.StatusConflict,
	}
)

// StartServer is start api server
func StartServer(ctx context.Context) {
	cfg := utils.Config()
	mux = bone.New()

	if cfg.DemoPage {
		mux.GetFunc("/", messengerHTMLHandler)
	}

	mux.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.GetFunc("/stats", stats_api.Handler)
	mux.GetFunc("/", indexHandler)
	mux.OptionsFunc("/*", optionsHandler)
	setAssetMux()
	setBlockUserMux()
	setDeviceMux()
	// setGuestMux()
	setMessageMux()
	setRoomMux()
	setRoomUserMux()
	setSettingMux()
	setUserMux()
	setUserRoleMux()

	if cfg.Profiling {
		setPprofMux()
	}

	if cfg.Storage.Provider == "awsS3" {
		setAssetAwsSnsMux()
	}

	mux.NotFoundFunc(notFoundHandler)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	errCh := make(chan error)
	go func() {
		errCh <- gracedown.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), mux)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[HTTP]", utils.AppName))
		gracedown.Close()
	case s := <-signalChan:
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			logger.Info(fmt.Sprintf("Stopping %s server[HTTP]", utils.AppName))
			gracedown.Close()
		}
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[HTTP]. %v", utils.AppName, err))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", fmt.Sprintf("%s [API Version]%s [Build Version]%s", utils.AppName, utils.APIVersion, utils.BuildVersion))
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

	rHeaders := make([]string, 0, len(r.Header))
	for k, v := range r.Header {
		rHeaders = append(rHeaders, k)
		if k == "Access-Control-Request-Headers" {
			rHeaders = append(rHeaders, strings.Join(v, ", "))
		}
	}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(rHeaders, ", "))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusNotFound, "", nil)
}

func commonHandler(fn http.HandlerFunc) http.HandlerFunc {
	return (colsHandler(
		jwtHandler(
			judgeAppClientHandler(
				func(w http.ResponseWriter, r *http.Request) {
					// log.Printf("url=%s\n", r.RequestURI)
					// domain := r.Host
					// referer := r.Header.Get("Referer")
					// if referer != "" {
					// 	url, err := url.Parse(referer)
					// 	if err != nil {
					// 		panic(err)
					// 	}
					// 	domain = url.Hostname()
					// }
					// log.Printf("domain=%s", domain)
					// for i, v := range r.Header {
					// 	log.Printf("%s=%s\n", i, v)
					// }
					fn(w, r)
				}))))
}

func colsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optionsHandler(w, r)
		fn(w, r)
	}
}

func jwtHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(utils.HeaderUserID)
		ctx := context.WithValue(r.Context(), utils.CtxUserID, userID)

		workspace := r.Header.Get(utils.HeaderWorkspace)
		ctx = context.WithValue(ctx, utils.CtxWorkspace, workspace)

		fn(w, r.WithContext(ctx))
	}
}

func updateLastAccessedHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)

		ctx := r.Context()
		userID := ctx.Value(utils.CtxUserID).(string)
		if userID == "" {
			return
		}

		go func() {
			user, err := datastore.Provider(ctx).SelectUser(userID)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			if user == nil {
				return
			}

			user.LastAccessed = time.Now().Unix()
			_, err = datastore.Provider(ctx).UpdateUser(user)
			if err != nil {
				logger.Error(err.Error())
			}
		}()
	}
}

func judgeAppClientHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := false
		clientID := r.Header.Get(utils.HeaderClientID)
		if clientID != "" {
			appCli, err := datastore.Provider(r.Context()).SelectLatestAppClientByClientID(clientID)
			if err != nil {
				respondErr(w, r, http.StatusInternalServerError, &model.ProblemDetail{
					Error: err,
				})
				return
			}
			if appCli != nil {
				isAppClient = true
			}
		}
		ctx := context.WithValue(r.Context(), utils.CtxIsAppClient, isAppClient)
		fn(w, r.WithContext(ctx))
	}
}

func adminAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := r.Context().Value(utils.CtxIsAppClient).(bool)
		if !isAppClient {
			respondErr(w, r, http.StatusUnauthorized, &model.ProblemDetail{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
			})
			return
		}
		fn(w, r)
	}
}

func selfResourceAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := r.Context().Value(utils.CtxIsAppClient).(bool)
		if isAppClient {
			fn(w, r)
			return
		}

		requestUserID := r.Header.Get(utils.HeaderUserID)
		resourceUserID := bone.GetValue(r, "userId")

		if requestUserID != "" && resourceUserID != "" {
			if requestUserID != resourceUserID {
				respondErr(w, r, http.StatusUnauthorized, &model.ProblemDetail{
					Message: fmt.Sprintf("Not your resource. Resource UserID is %s, but request UserID is %s.", resourceUserID, requestUserID),
					Status:  http.StatusUnauthorized,
				})
				return
			}
		}
		fn(w, r)
	}
}

func contactsAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := r.Context().Value(utils.CtxIsAppClient).(bool)
		if isAppClient {
			fn(w, r)
			return
		}

		requestUserID := r.Header.Get(utils.HeaderUserID)
		resourceUserID := bone.GetValue(r, "userId")
		pd := service.ContactsAuthz(r.Context(), requestUserID, resourceUserID)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		fn(w, r)
	}
}

func roomMemberAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := r.Context().Value(utils.CtxIsAppClient).(bool)
		if isAppClient {
			fn(w, r)
			return
		}

		roomID := bone.GetValue(r, "roomId")
		userID := r.Header.Get(utils.HeaderUserID)

		pd := service.RoomAuthz(r.Context(), roomID, userID)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		fn(w, r)
	}
}

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	bodySize, _ := buf.ReadFrom(r.Body)
	if bodySize == 0 {
		return nil
	}
	return json.NewDecoder(buf).Decode(v)
}

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request, status int, contentType string, data interface{}) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(status)
	for _, v := range noBodyStatusCodes {
		if status == v {
			data = nil
		}
	}
	if data != nil {
		encodeBody(w, r, data)
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, status int, pd *model.ProblemDetail) {
	if pd.Error != nil {
		logger.Error(pd.Error.Error())
		if utils.Config().EnableDeveloperMessage {
			pd.DeveloperMessage = pd.Error.Error()
		}
	}
	respond(w, r, status, "application/json", pd)
}

func respondJSONDecodeError(w http.ResponseWriter, r *http.Request, title string) {
	respondErr(w, r, http.StatusBadRequest, &model.ProblemDetail{
		Message: fmt.Sprintf("Json parse error. (%s)", title),
		Status:  http.StatusBadRequest,
	})
}

func setLastModified(w http.ResponseWriter, timestamp int64) {
	l, _ := time.LoadLocation("Etc/GMT")
	lm := time.Unix(timestamp, 0).In(l).Format(http.TimeFormat)
	w.Header().Set("Last-Modified", lm)
}

func setPagingParams(params url.Values) (int32, int32, string, *model.ProblemDetail) {
	limit := int32(10)
	offset := int32(0)
	order := "ASC"
	if limitArray, ok := params["limit"]; ok {
		limitInt, err := strconv.Atoi(limitArray[0])
		if err != nil {
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "limit",
						Reason: "limit is incorrect.",
					},
				},
			}
		}
		limit = int32(limitInt)
	}
	if offsetArray, ok := params["offset"]; ok {
		offsetInt, err := strconv.Atoi(offsetArray[0])
		if err != nil {
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "offset",
						Reason: "offset is incorrect.",
					},
				},
			}
		}
		offset = int32(offsetInt)
	}
	if orderArray, ok := params["order"]; ok {
		order := orderArray[0]
		allowedOrders := []string{
			"DESC",
			"desc",
			"ASC",
			"asc",
		}
		if utils.SearchStringValueInSlice(allowedOrders, order) {
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
						Name:   "order",
						Reason: "order is incorrect.",
					},
				},
			}
		}
	}
	return limit, offset, order, nil
}
