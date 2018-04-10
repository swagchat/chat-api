// http handler
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
	"github.com/swagchat/chat-api/utils"
)

type key int

const (
	jwtSub       = "X-Sub"
	jwtRealm     = "X-Realm"
	ctxDsCfg key = iota
	ctxJwt
)

var (
	Mux            *bone.Mux
	Context        context.Context
	allowedMethods []string = []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"PATCH",
		"DELETE",
	}
	NoBodyStatusCodes []int = []int{
		http.StatusNotFound,
		http.StatusConflict,
	}
)

func StartServer(ctx context.Context) {
	cfg := utils.Config()

	Mux = bone.New()
	if cfg.DemoPage {
		Mux.GetFunc("/", messengerHTMLHandler)
	}
	Mux.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	Mux.GetFunc("/stats", stats_api.Handler)
	Mux.GetFunc("/", indexHandler)
	Mux.OptionsFunc("/*", optionsHandler)
	SetAssetMux()
	SetBlockUserMux()
	SetContactMux()
	SetDeviceMux()
	SetMessageMux()
	SetRoomMux()
	SetRoomUserMux()
	SetSettingMux()
	SetUserMux()

	if cfg.Profiling {
		SetPprofMux()
	}

	if cfg.Storage.Provider == "awsS3" {
		SetAssetAwsSnsMux()
	}

	Mux.NotFoundFunc(notFoundHandler)

	go run(ctx)

	c := utils.Config()
	sb := utils.NewStringBuilder()
	cfgStr := sb.PrintStruct("config", c)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "handler",
		Message: fmt.Sprintf("%s start", utils.AppName),
		Config:  cfgStr,
	})

	err := gracedown.ListenAndServe(utils.AppendStrings(":", cfg.HttpPort), Mux)
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

func colsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optionsHandler(w, r)
		fn(w, r)
	}
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	rHeaders := make([]string, 0, len(r.Header))
	for k, v := range r.Header {
		rHeaders = append(rHeaders, k)
		if k == "Access-Control-Request-Headers" {
			rHeaders = append(rHeaders, strings.Join(v, ", "))
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(rHeaders, ", "))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusNotFound, "", nil)
}

func jwtHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt := &models.JWT{
			Sub: r.Header.Get(jwtSub),
		}

		ctx := context.WithValue(r.Context(), ctxJwt, jwt)
		fn(w, r.WithContext(ctx))
	}
}

func datastoreHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := utils.Config()
		var dsCfg *utils.Datastore

		switch utils.Config().Datastore.Provider {
		case "sqlite":
			dsCfg = &utils.Datastore{
				Provider:        cfg.Datastore.Provider,
				TableNamePrefix: cfg.Datastore.TableNamePrefix,
			}
			dsCfg.SQLite.Path = cfg.Datastore.SQLite.Path
		case "mysql", "gcSql":
			dsCfg = &utils.Datastore{
				Provider:          cfg.Datastore.Provider,
				User:              cfg.Datastore.User,
				Password:          cfg.Datastore.Password,
				Database:          cfg.Datastore.Database,
				TableNamePrefix:   cfg.Datastore.TableNamePrefix,
				Master:            cfg.Datastore.Master,
				Replicas:          cfg.Datastore.Replicas,
				MaxIdleConnection: cfg.Datastore.MaxIdleConnection,
				MaxOpenConnection: cfg.Datastore.MaxOpenConnection,
			}
		}
		if cfg.Datastore.Dynamic {
			dsCfg.Database = r.Header.Get(jwtRealm)
		}

		if dsCfg.Database == "" {
			respondErr(w, r, 400, &models.ProblemDetail{
				Title: "No database",
			})
			return
		}

		datastore.Provider(dsCfg)
		ctx := context.WithValue(r.Context(), ctxDsCfg, dsCfg)
		fn(w, r.WithContext(ctx))
	}
}

// func aclHandler(fn http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		role := "guest"

// 		apiKey := r.Header.Get(utils.HeaderAPIKey)
// 		apiSecret := r.Header.Get(utils.HeaderAPISecret)
// 		if apiKey != "" && apiSecret != "" {
// 			api, err := datastore.Provider().SelectLatestApi("admin")
// 			if err != nil {
// 				// TODO error
// 			}
// 			if api != nil {
// 				if apiKey == api.Key && apiSecret == api.Secret {
// 					role = "admin"
// 				}
// 			}
// 		}

// 		if role != "admin" {
// 			authorization := r.Header.Get("Authorization")
// 			token := strings.Replace(authorization, "Bearer ", "", 1)
// 			userId := r.Header.Get(utils.HeaderUserId)
// 			if token != "" && userId != "" {
// 				user, err := datastore.Provider().SelectUserByUserIdAndAccessToken(userId, token)
// 				if err != nil {
// 					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
// 						Message:    err.Error(),
// 						Stacktrace: fmt.Sprintf("%v\n", err),
// 					})
// 				}
// 				if user != nil {
// 					role = "user"
// 				}
// 			}
// 		}
// 		ctx := context.WithValue(r.Context(), "role", role)
// 		fn(w, r.WithContext(ctx))
// 	}
// }

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
	for _, v := range NoBodyStatusCodes {
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

func respondJsonDecodeError(w http.ResponseWriter, r *http.Request, title string) {
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
