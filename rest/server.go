package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
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

// Run runs start REST API server
func Run(ctx context.Context) {
	cfg := config.Config()
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
	setMessageMux()
	setRoomMux()
	setRoomUserMux()
	setSettingMux()
	setUserMux()
	setUserRoleMux()

	if cfg.Profiling {
		mux.GetFunc("/debug/vars", metricsHandler)
	}

	if cfg.Storage.Provider == "awsS3" {
		setAssetAwsSnsMux()
	}

	mux.NotFoundFunc(notFoundHandler)

	logger.Info(fmt.Sprintf("Starting %s server[REST] on listen tcp :%s", config.AppName, cfg.HTTPPort))
	errCh := make(chan error)
	go func() {
		errCh <- gracedown.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), mux)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[REST]", config.AppName))
		datastore.Provider(ctx).Close()
		gracedown.Close()
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[REST]. %v", config.AppName, err))
		datastore.Provider(ctx).Close()
		gracedown.Close()
	}
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", fmt.Sprintf("%s [API Version]%s [Build Version]%s", config.AppName, config.APIVersion, config.BuildVersion))
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
		traceHandler(
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
						defer r.Body.Close()
						fn(w, r)
					})))))
}

func colsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optionsHandler(w, r)
		fn(w, r)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "%q: %q", key, str)
		} else {
			fmt.Fprintf(w, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func traceHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Provider(r.Context()).StartTransaction(
			fmt.Sprintf("%s:%s", r.Method, r.RequestURI), "REST",
			tracer.StartTransactionOptionWithHTTPRequest(r),
		)
		defer tracer.Provider(ctx).CloseTransaction()

		sw := &customResponseWriter{ResponseWriter: w}
		fn(sw, r.WithContext(ctx))

		userID := ctx.Value(config.CtxUserID)
		if userID != nil {
			tracer.Provider(ctx).SetTag(span, "userId", userID)
		}
		clientID := ctx.Value(config.CtxClientID)
		if clientID != nil {
			tracer.Provider(ctx).SetTag(span, "clientId", clientID)
		}
		tracer.Provider(ctx).SetHTTPStatusCode(span, sw.status)
		tracer.Provider(ctx).SetTag(span, "http.method", r.Method)
		tracer.Provider(ctx).SetTag(span, "http.content_length", sw.length)
		tracer.Provider(ctx).SetTag(span, "http.referer", r.Referer())
	}
}

func jwtHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(config.HeaderUserID)
		ctx := context.WithValue(r.Context(), config.CtxUserID, userID)

		workspace := r.Header.Get(config.HeaderWorkspace)
		ctx = context.WithValue(ctx, config.CtxWorkspace, workspace)

		fn(w, r.WithContext(ctx))
	}
}

func updateLastAccessedHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)

		ctx := r.Context()
		userID := ctx.Value(config.CtxUserID).(string)
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

			user.LastAccessedTimestamp = time.Now().Unix()
			err = datastore.Provider(ctx).UpdateUser(user)
			if err != nil {
				logger.Error(err.Error())
			}
		}()
	}
}

func judgeAppClientHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get(config.HeaderClientID)
		if clientID != "" {
			appCli, err := datastore.Provider(r.Context()).SelectLatestAppClient(
				datastore.SelectAppClientOptionFilterByClientID(clientID),
			)
			if err != nil {
				errRes := model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
				respondError(w, r, errRes)
				return
			}
			if appCli != nil {
				clientID = appCli.ClientID
			}
		}
		ctx := context.WithValue(r.Context(), config.CtxClientID, clientID)
		fn(w, r.WithContext(ctx))
	}
}

func adminAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(config.CtxClientID)
		if clientID == "" {
			errRes := model.NewErrorResponse("Unauthorized", http.StatusUnauthorized)
			respondError(w, r, errRes)
			return
		}
		fn(w, r)
	}
}

func selfResourceAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(config.CtxClientID)
		if clientID != "" {
			fn(w, r)
			return
		}

		requestUserID := r.Header.Get(config.HeaderUserID)
		resourceUserID := bone.GetValue(r, "userId")

		if (requestUserID == "" && resourceUserID == "") || (requestUserID != resourceUserID) {
			errRes := model.NewErrorResponse(fmt.Sprintf("Not your resource. Resource UserID is %s, but request UserID is %s.", resourceUserID, requestUserID), http.StatusUnauthorized)
			respondError(w, r, errRes)
			return
		}

		fn(w, r)
	}
}

func contactsAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(config.CtxClientID)
		if clientID != "" {
			fn(w, r)
			return
		}

		requestUserID := r.Header.Get(config.HeaderUserID)
		resourceUserID := bone.GetValue(r, "userId")
		errRes := service.ContactsAuthz(r.Context(), requestUserID, resourceUserID)
		if errRes != nil {
			respondError(w, r, errRes)
			return
		}
		fn(w, r)
	}
}

func roomMemberAuthzHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(config.CtxClientID)
		if clientID != "" {
			fn(w, r)
			return
		}

		roomID := bone.GetValue(r, "roomId")
		userID := r.Header.Get(config.HeaderUserID)

		errRes := service.RoomAuthz(r.Context(), roomID, userID)
		if errRes != nil {
			respondError(w, r, errRes)
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

func respondError(w http.ResponseWriter, r *http.Request, errRes *model.ErrorResponse) {
	if errRes.Error != nil {
		if config.Config().EnableDeveloperMessage {
			errRes.DeveloperMessage = errRes.Error.Error()
		}
	}
	respond(w, r, errRes.Status, "application/json", errRes)
}

func respondJSONDecodeError(w http.ResponseWriter, r *http.Request, title string) {
	errRes := model.NewErrorResponse(fmt.Sprintf("Json parse error. (%s)", title), http.StatusBadRequest)
	respondError(w, r, errRes)
}

func setLastModified(w http.ResponseWriter, timestamp int64) {
	l, _ := time.LoadLocation("Etc/GMT")
	lm := time.Unix(timestamp, 0).In(l).Format(http.TimeFormat)
	w.Header().Set("Last-Modified", lm)
}

func setPagingParams(params url.Values) (int32, int32, int64, int64, []*scpb.OrderInfo, *model.ErrorResponse) {
	limit := int32(10)
	offset := int32(0)
	limitTimestamp := int64(0)
	offsetTimestamp := int64(0)
	var orders []*scpb.OrderInfo

	if limitArray, ok := params["limit"]; ok {
		limitInt, err := strconv.Atoi(limitArray[0])
		if err != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "limit",
					Reason: "limit is incorrect.",
				},
			}
			return limit, offset, limitTimestamp, offsetTimestamp, orders, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		}
		limit = int32(limitInt)
	}
	if offsetArray, ok := params["offset"]; ok {
		offsetInt, err := strconv.Atoi(offsetArray[0])
		if err != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "offset",
					Reason: "offset is incorrect.",
				},
			}
			return limit, offset, limitTimestamp, offsetTimestamp, orders, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		}
		offset = int32(offsetInt)
	}
	if limitTimestampArray, ok := params["limitTimestamp"]; ok {
		limitTimestampInt, err := strconv.Atoi(limitTimestampArray[0])
		if err != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "limitTimestamp",
					Reason: "limitTimestamp is incorrect.",
				},
			}
			return limit, offset, limitTimestamp, offsetTimestamp, orders, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		}
		limitTimestamp = int64(limitTimestampInt)
	}
	if offsetTimestampArray, ok := params["offsetTimestamp"]; ok {
		offsetTimestampInt, err := strconv.Atoi(offsetTimestampArray[0])
		if err != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "offsetTimestamp",
					Reason: "offsetTimestamp is incorrect.",
				},
			}
			return limit, offset, limitTimestamp, offsetTimestamp, orders, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		}
		offsetTimestamp = int64(offsetTimestampInt)
	}
	if orderArray, ok := params["order"]; ok {
		orderString := orderArray[0] // ex) field1+desc,field2+asc

		orderPairs := strings.Split(orderString, ",")
		orders := make([]*scpb.OrderInfo, len(orderPairs))
		for _, orderPair := range orderPairs {
			order := strings.Split(orderPair, " ")
			if len(order) != 2 {
				continue
			}
			if orderInt32, ok := scpb.Order_value[order[1]]; ok {
				orderInfo := &scpb.OrderInfo{
					Field: order[0],
					Order: scpb.Order(orderInt32),
				}
				orders = append(orders, orderInfo)
			}
		}
	}

	return limit, offset, limitTimestamp, offsetTimestamp, orders, nil
}
