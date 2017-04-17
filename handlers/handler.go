// http handler
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
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

func StartServer() {
	Mux = bone.New().Prefix(utils.AppendStrings("/", utils.API_VERSION))
	Mux.GetFunc("", indexHandler)
	Mux.GetFunc("/", indexHandler)
	Mux.GetFunc("/sleep", sleepHandler2)
	Mux.GetFunc("/stats", stats_api.Handler)
	Mux.OptionsFunc("/*", optionsHandler)
	SetUserMux()
	SetRoomMux()
	SetRoomUserMux()
	SetMessageMux()
	SetAssetMux()
	SetDeviceMux()
	if utils.Cfg.ApiServer.Storage == "awsS3" {
		SetAssetAwsSnsMux()
	}
	Mux.NotFoundFunc(notFoundHandler)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				utils.AppLogger.Info("",
					zap.String("msg", "Swagchat API Shutdown start!"),
					zap.String("signal", s.String()),
				)
				gracedown.Close()
			}
		}
	}()

	utils.AppLogger.Info("",
		zap.String("msg", "Swagchat API Start!"),
		zap.String("port", utils.Cfg.ApiServer.Port),
	)
	if err := gracedown.ListenAndServe(utils.AppendStrings(":", utils.Cfg.ApiServer.Port), Mux); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	utils.AppLogger.Info("",
		zap.String("msg", "Swagchat API Shutdown finish!"),
	)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusCreated, "text/plain", utils.AppendStrings("Swagchat API version ", utils.API_VERSION))
}

func sleepHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(5 * time.Second)
			//xxx(r.Context())
		}()
	}
	wg.Wait()
	respond(w, r, http.StatusCreated, "text/plain", utils.AppendStrings("Swagchat API version ", utils.API_VERSION))
}

func sleepHandler2(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//ctx := context.Background()

	//sigCh := make(chan os.Signal, 1)
	//defer close(sigCh)

	//signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	//go func() {
	//	<-sigCh
	//	cancel()
	//}()

	d := utils.NewDispatcher(5)
	for i := 0; i < 20; i++ {
		d.Work(ctx, func(ctx context.Context) {
			log.Printf("start procing %d", i)

			abcChan := make(chan *abc, 1)
			abcObj := xxx()
			abcChan <- abcObj

			select {
			case <-ctx.Done():
				log.Printf("cancel work func")
				return
			case <-abcChan:
				log.Printf("done procing")
				return
			}
		})
	}

	d.Wait()
	respond(w, r, http.StatusCreated, "text/plain", utils.AppendStrings("Swagchat API version ", utils.API_VERSION))
}

type abc struct {
}

func xxx() *abc {
	time.Sleep(3 * time.Second)
	return &abc{}
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
	utils.AppLogger.Error("",
		zap.String("msg", "404 Not Found"),
	)
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
	problemDetailBytes, _ := json.Marshal(problemDetail)
	utils.AppLogger.Error("",
		zap.String("problemDetail", string(problemDetailBytes)),
		zap.String("err", fmt.Sprintf("%+v", problemDetail.Error)),
	)
	respond(w, r, status, "application/json", problemDetail)
}

func respondJsonDecodeError(w http.ResponseWriter, r *http.Request, title string) {
	respondErr(w, r, http.StatusBadRequest, &models.ProblemDetail{
		Title:     utils.AppendStrings("Json parse error. (", title, ")"),
		Status:    http.StatusBadRequest,
		ErrorName: models.ERROR_NAME_INVALID_JSON,
	})
}
