package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

var (
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

func StartServer() {
	mux := bone.New()
	mux.GetFunc("", indexHandler)
	mux.GetFunc("/", indexHandler)
	mux.GetFunc("/stats", stats_api.Handler)
	mux.GetFunc(utils.AppendStrings("/", utils.API_VERSION), websocketHandler)
	mux.GetFunc(utils.AppendStrings("/", utils.API_VERSION, "/"), websocketHandler)
	mux.PostFunc(utils.AppendStrings("/", utils.API_VERSION, "/message"), messageHandler)
	mux.OptionsFunc(utils.AppendStrings("/", utils.API_VERSION, "/*"), optionsHandler)
	mux.NotFoundFunc(notFoundHandler)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Println("msg", "Swagchat Realtime Shutdown start!")
				messaging.Unsubscribe()
				gracedown.Close()
			}
		}
	}()
	aa := `
███████╗██╗    ██╗ █████╗  ██████╗  ██████╗██╗  ██╗ █████╗ ████████╗    ██████╗ ████████╗███╗   ███╗     █████╗ ██████╗ ██╗
██╔════╝██║    ██║██╔══██╗██╔════╝ ██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██╔══██╗╚══██╔══╝████╗ ████║    ██╔══██╗██╔══██╗██║
███████╗██║ █╗ ██║███████║██║  ███╗██║     ███████║███████║   ██║       ██████╔╝   ██║   ██╔████╔██║    ███████║██████╔╝██║
╚════██║██║███╗██║██╔══██║██║   ██║██║     ██╔══██║██╔══██║   ██║       ██╔══██╗   ██║   ██║╚██╔╝██║    ██╔══██║██╔═══╝ ██║
███████║╚███╔███╔╝██║  ██║╚██████╔╝╚██████╗██║  ██║██║  ██║   ██║       ██║  ██║   ██║   ██║ ╚═╝ ██║    ██║  ██║██║     ██║
╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝    ╚═╝  ╚═╝╚═╝     ╚═╝

`
	fmt.Println(aa)
	fmt.Printf("port %s\n", utils.Realtime.Port)
	log.Println("Swagchat realtime server start!")
	if err := gracedown.ListenAndServe(utils.AppendStrings(":", utils.Realtime.Port), mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", utils.AppendStrings("swagchat RTM API version ", utils.API_VERSION))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- messageHandler ---")
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	services.Srv.Broadcast <- message
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &services.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	params, _ := url.ParseQuery(r.URL.RawQuery)
	if userIdSlice, ok := params["userId"]; ok {
		c.UserId = userIdSlice[0]
	}

	services.Srv.Connection.AddClient(c)
	go c.WritePump()
	go c.ReadPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("options")
	log.Println(r.Header.Get("Access-Control-Request-Headers"))

	origin := r.Header.Get("Origin")
	if origin != "" {
		log.Println(origin)
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
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
