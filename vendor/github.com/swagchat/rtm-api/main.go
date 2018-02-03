// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/swagchat/rtm-api/handlers"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

func main() {
	var (
		port           string
		nsqlookupdHost string
		nsqlookupdPort string
		nsqdHost       string
		nsqdPort       string
		queTopic       string
		queChannel     string
	)
	log.SetFlags(log.Lshortfile)

	var v string
	if v = os.Getenv("PORT"); v != "" {
		port = v
	} else {
		flag.StringVar(&port, "port", "9100", "service port")
	}
	flag.StringVar(&nsqlookupdHost, "nsqlookupdHost", "", "Host name of nsqlookupd")
	flag.StringVar(&nsqlookupdPort, "nsqlookupdPort", "4161", "Port no of nsqlookupd")
	flag.StringVar(&nsqdHost, "nsqdHost", "", "Host name of nsqd")
	flag.StringVar(&nsqdPort, "nsqdPort", "4151", "Port no of nsqd")
	flag.StringVar(&queTopic, "queTopic", "websocket", "Topic name")
	flag.StringVar(&queChannel, "queChannel", "", "Channel name. If it's not set, channel is hostname.")
	flag.Parse()

	utils.Realtime.Port = port
	utils.Que.NsqlookupdHost = nsqlookupdHost
	utils.Que.NsqlookupdPort = nsqlookupdPort
	utils.Que.NsqdHost = nsqdHost
	utils.Que.NsqdPort = nsqdPort
	utils.Que.Topic = queTopic
	utils.Que.Channel = queChannel

	messaging.Subscribe()
	go services.Srv.Run()
	handlers.StartServer()
}
