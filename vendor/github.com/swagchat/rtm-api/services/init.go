package services

func init() {
	connection := Connection{
		clients: make(map[*Client]bool),
		users:   make(map[string]UserClients),
		rooms:   make(map[string]RoomClients),
	}

	Srv = Server{
		Connection: connection,
		Broadcast:  make(chan []byte),
		Register:   make(chan *RcvData),
		Unregister: make(chan *RcvData),
		Close:      make(chan *Client),
	}
}
