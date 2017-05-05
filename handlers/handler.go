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

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/shogo82148/go-gracedown"
)

var Mux *bone.Mux
var Context context.Context

func StartServer(ctx context.Context) {
	Mux = bone.New().Prefix(utils.AppendStrings("/", utils.API_VERSION))
	Mux.GetFunc("", indexHandler)
	Mux.GetFunc("/", indexHandler)
	Mux.GetFunc("/stats", stats_api.Handler)
	Mux.OptionsFunc("/*", optionsHandler)
	SetUserMux()
	SetRoomMux()
	SetRoomUserMux()
	SetMessageMux()
	SetAssetMux()
	SetDeviceMux()
	if utils.Cfg.Storage.Provider == "awsS3" {
		SetAssetAwsSnsMux()
	}
	Mux.NotFoundFunc(notFoundHandler)

	go run(ctx)

	utils.AppLogger.Info("",
		zap.String("msg", "Swagchat API Start!"),
		zap.String("port", utils.Cfg.Port),
	)
	if err := gracedown.ListenAndServe(utils.AppendStrings(":", utils.Cfg.Port), Mux); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	utils.AppLogger.Info("",
		zap.String("msg", "Swagchat API Shutdown finish!"),
	)
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
				utils.AppLogger.Info("",
					zap.String("msg", "Swagchat API Shutdown start!"),
					zap.String("signal", s.String()),
				)
				gracedown.Close()
			}
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", utils.AppendStrings("Swagchat API version ", utils.API_VERSION))
}

func ColsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	utils.AppLogger.Info("",
		zap.String("msg", "call optionHandler"),
		zap.String("origin", origin),
	)
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusNotFound, "", nil)
}

var allowedMethods []string = []string{
	"POST",
	"GET",
	"OPTIONS",
	"PUT",
	"PATCH",
	"DELETE",
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
	if status == http.StatusNotFound {
		data = nil
	}
	if data != nil {
		encodeBody(w, r, data)
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, status int, problemDetail *models.ProblemDetail) {
	if !utils.Cfg.ErrorLogging {
		problemDetailBytes, _ := json.Marshal(problemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(problemDetailBytes)),
			zap.String("err", fmt.Sprintf("%+v", problemDetail.Error)),
		)
	}
	respond(w, r, status, "application/json", problemDetail)
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
