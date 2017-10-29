package main

import (

	"image"
)


type ClientManager struct {
	clients      map[*Client]bool
	picture      chan image.Image
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	picBroadcast chan map[*Client][]byte


}

func newManager() *ClientManager {
	return &ClientManager{
		picture:      make(chan image.Image),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		}
}


func (manager *ClientManager) send(message []byte) {
	for conn := range manager.clients {
		conn.send <- message

	}
}

func (manager *ClientManager) start() {


	for {
		select {
		case cliReg := <-manager.register:
			manager.clients[cliReg] = true
			//jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			//manager.send(jsonMessage, cliReg)
		case clieUnReg := <-manager.unregister:
			if _, ok := manager.clients[clieUnReg]; ok {
				close(clieUnReg.send)
				delete(manager.clients, clieUnReg)
				//		jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				//		manager.send(jsonMessage, cliReg)
			}
		case pic := <-manager.picture:
			for client := range manager.clients {
				client.prossPic <- pic
			}


		}

	}
}