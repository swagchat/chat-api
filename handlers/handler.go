package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
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
	configValidate()

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
	setGuestMux()
	setMessageMux()
	setRoomMux()
	setRoomUserMux()
	setSettingMux()
	setUserMux()

	if cfg.Profiling {
		setPprofMux()
	}

	if cfg.Storage.Provider == "awsS3" {
		setAssetAwsSnsMux()
	}

	mux.NotFoundFunc(notFoundHandler)

	go run(ctx)

	sb := utils.NewStringBuilder()
	cfgStr := sb.PrintStruct("config", cfg)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "handler",
		Message: fmt.Sprintf("%s start", utils.AppName),
		Config:  cfgStr,
	})

	err := gracedown.ListenAndServe(utils.AppendStrings(":", cfg.HTTPPort), mux)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}

	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "handler",
		Message: "shut down complete",
	})
}

func run(ctx context.Context) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-ctx.Done():
			gracedown.Close()
		case s := <-signalChan:
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:    "handler",
					Message: fmt.Sprintf("%s graceful down by signal[%s]", utils.AppName, s.String()),
				})
				gracedown.Close()
			}
		}
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
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: err,
				})
				return
			}

			if user == nil {
				return
			}

			user.LastAccessed = time.Now().Unix()
			_, err = datastore.Provider(ctx).UpdateUser(user)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: err,
				})
			}
		}()
	}
}

func judgeAppClientHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAppClient := false
		username := r.Header.Get(utils.HeaderUsername)
		if username != "" {
			api, err := datastore.Provider(r.Context()).SelectLatestAppClientByClientID(username)
			if err != nil {
				respondErr(w, r, http.StatusInternalServerError, &models.ProblemDetail{
					Error: err,
				})
				return
			}
			if api != nil {
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
			respondErr(w, r, http.StatusUnauthorized, &models.ProblemDetail{
				Status: http.StatusUnauthorized,
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
				respondErr(w, r, http.StatusUnauthorized, &models.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Not your resource",
					Detail: fmt.Sprintf("Resource UserID is %s, but request UserID is %s.", resourceUserID, requestUserID),
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
		pd := services.ContactsAuthz(r.Context(), requestUserID, resourceUserID)
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

		pd := services.RoomAuthz(r.Context(), roomID, userID)
		if pd != nil {
			respondErr(w, r, pd.Status, pd)
			return
		}
		fn(w, r)
	}
}

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	bufbody := new(bytes.Buffer)
	bodySize, _ := bufbody.ReadFrom(r.Body)
	if bodySize == 0 {
		return nil
	}
	return json.NewDecoder(bufbody).Decode(v)
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

func respondErr(w http.ResponseWriter, r *http.Request, status int, pd *models.ProblemDetail) {
	if pd.Error != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: pd.Error,
		})
	}
	respond(w, r, status, "application/json", pd)
}

func respondJSONDecodeError(w http.ResponseWriter, r *http.Request, title string) {
	respondErr(w, r, http.StatusBadRequest, &models.ProblemDetail{
		Title:     utils.AppendStrings("Json parse error. (", title, ")"),
		Status:    http.StatusBadRequest,
		ErrorName: models.ERROR_NAME_INVALID_JSON,
	})
}

func setLastModified(w http.ResponseWriter, timestamp int64) {
	l, _ := time.LoadLocation("Etc/GMT")
	lm := time.Unix(timestamp, 0).In(l).Format(http.TimeFormat)
	w.Header().Set("Last-Modified", lm)
}
