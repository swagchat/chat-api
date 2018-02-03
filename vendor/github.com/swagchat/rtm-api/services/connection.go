package services

import (
	"log"
	"os"
)

type Connection struct {
	clients map[*Client]bool
	users   map[string]UserClients // index is userId
	rooms   map[string]RoomClients // index is roomId
}

type UserClients struct {
	clients map[*Client]bool
}

type RoomClients struct {
	roomUsers map[string]RoomUserClients // index is userId
}

type RoomUserClients struct {
	events map[string]EventClients // index is eventName
}

type EventClients struct {
	clients map[*Client]bool
}

func (con *Connection) AddClient(c *Client) {
	if c == nil {
		return
	}

	hostname, _ := os.Hostname()
	log.Printf("[WS-INFO][%s] ADD CLIENT %p", hostname, c)

	con.clients[c] = true

	var userClients UserClients
	if _, ok := con.users[c.UserId]; ok {
		con.users[c.UserId].clients[c] = true
	} else {
		userClients = UserClients{
			clients: make(map[*Client]bool),
		}
		userClients.clients[c] = true
		con.users[c.UserId] = userClients
	}
}

func (con *Connection) Info() {
	hostname, _ := os.Hostname()
	log.Printf("[WS-INFO][%s] All Clients %d", hostname, len(con.clients))
	for userId, _ := range con.users {
		log.Printf("[WS-INFO][%s] UserId[%s] %d", hostname, userId, len(con.users[userId].clients))
	}
	for roomId, _ := range con.rooms {
		log.Printf("[WS-INFO][%s] RoomId[%s] %d", hostname, roomId, len(con.rooms[roomId].roomUsers))
		for userId, _ := range con.rooms[roomId].roomUsers {
			log.Printf("[WS-INFO][%s] RoomId[%s][%s] %d", hostname, roomId, userId, len(con.rooms[roomId].roomUsers[userId].events))
			for eventName, _ := range con.rooms[roomId].roomUsers[userId].events {
				log.Printf("[WS-INFO][%s] RoomId[%s][%s][%s] %d", hostname, roomId, userId, eventName, len(con.rooms[roomId].roomUsers[userId].events[eventName].clients))
			}
		}
	}
}

func (con *Connection) RemoveClient(c *Client) {
	if c == nil {
		return
	}

	delete(con.clients, c)

	for client, _ := range con.users[c.UserId].clients {
		if client == c {
			delete(con.users[c.UserId].clients, c)
		}
	}
	if len(con.users[c.UserId].clients) == 0 {
		delete(con.users, c.UserId)
	}

	for roomId, roomClients := range con.rooms {
		for userId, roomUserClients := range roomClients.roomUsers {
			for eventName, eventClients := range roomUserClients.events {
				for client, _ := range eventClients.clients {
					if client == c {
						delete(eventClients.clients, c)
					}
				}
				if len(roomUserClients.events[eventName].clients) == 0 {
					delete(roomUserClients.events, eventName)
				}
			}
			if len(roomClients.roomUsers[userId].events) == 0 {
				delete(roomClients.roomUsers, userId)
			}
		}
		if len(con.rooms[roomId].roomUsers) == 0 {
			delete(con.rooms, roomId)
		}
	}

}

func (con *Connection) AddEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId]; !ok {
		rc := RoomClients{
			roomUsers: make(map[string]RoomUserClients),
		}
		ruc := RoomUserClients{
			events: make(map[string]EventClients),
		}
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.events[eventName] = ec
		rc.roomUsers[userId] = ruc
		con.rooms[roomId] = rc
	} else if _, ok := con.rooms[roomId].roomUsers[userId]; !ok {
		ruc := RoomUserClients{
			events: make(map[string]EventClients),
		}
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.events[eventName] = ec
		con.rooms[roomId].roomUsers[userId] = ruc
	} else if _, ok := con.rooms[roomId].roomUsers[userId].events[eventName]; !ok {
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		con.rooms[roomId].roomUsers[userId].events[eventName] = ec
	} else {
		con.rooms[roomId].roomUsers[userId].events[eventName].clients[c] = true
	}
}

func (con *Connection) RemoveEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId].roomUsers[userId].events[eventName].clients[c]; ok {
		delete(con.rooms[roomId].roomUsers[userId].events[eventName].clients, c)
	}
	if len(con.rooms[roomId].roomUsers[userId].events[eventName].clients) == 0 {
		delete(con.rooms[roomId].roomUsers[userId].events, eventName)
	}
	if len(con.rooms[roomId].roomUsers[userId].events) == 0 {
		delete(con.rooms[roomId].roomUsers, userId)
	}
	if len(con.rooms[roomId].roomUsers) == 0 {
		delete(con.rooms, roomId)
	}
}
