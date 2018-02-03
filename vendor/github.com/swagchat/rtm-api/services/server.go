package services

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/swagchat/rtm-api/models"
)

var Srv Server

type Server struct {
	Connection
	Broadcast  chan []byte
	Register   chan *RcvData
	Unregister chan *RcvData
	Close      chan *Client
}

func (s *Server) Run() {
	hostname, _ := os.Hostname()

	for {
		var infoInterval <-chan time.Time
		infoInterval = time.After(5 * time.Second)

		select {
		case <-infoInterval:
			// Logging connection information.
			s.Connection.Info()

		case rcvData := <-s.Register:
			// Register event
			log.Printf("[WS-INFO][%s] REGISTER [%s][%s][%s] %p", hostname, rcvData.RoomId, rcvData.UserId, rcvData.EventName, rcvData.Client)
			s.clients[rcvData.Client] = true
			s.Connection.AddEvent(rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client)

		case rcvData := <-s.Unregister:
			// Unregister event
			log.Printf("[WS-INFO][%s] UNREGISTER [%s][%s][%s] %p", hostname, rcvData.RoomId, rcvData.UserId, rcvData.EventName, rcvData.Client)
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client)

		case c := <-s.Close:
			// Socket close
			log.Printf("[WS-INFO][%s] CLOSE UserId[%s]", hostname, c.UserId)
			s.Connection.RemoveClient(c)
			close(c.Send)

		case message := <-s.Broadcast:
			// Broadcast message
			log.Printf("[WS-INFO][%s] BROADCAST [%s]", hostname, string(message))
			s.broadcast(message)
		}
	}
}

func (s *Server) broadcast(message []byte) {
	var messageMap models.Message
	json.Unmarshal(message, &messageMap)
	if messageMap.Type == "text" {
		var payloadText models.PayloadText
		json.Unmarshal(messageMap.Payload, &payloadText)
	}

	for _, roomUser := range s.rooms[messageMap.RoomId].roomUsers {
		for conn, _ := range roomUser.events[messageMap.EventName].clients {
			log.Printf("------> %p", conn)
			select {
			case conn.Send <- message:
			default:
				close(conn.Send)
				delete(Srv.clients, conn)
			}
		}
	}
}
