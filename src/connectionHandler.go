package main

import (
	"image"
)

type ClientHandler struct {
	clients      map[*Client]bool
	picture      chan image.Image
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	picBroadcast chan map[*Client][]byte
	serverConf *serverConfiguration
	charConf      map[string] *chartConfig

}

func newClientHandler() *ClientHandler {

	serverConf := new (serverConfiguration)
	serverConf.loadConfig()
	return &ClientHandler{
		picture:    make(chan image.Image),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		serverConf: serverConf,
		charConf:	make(map[string] *chartConfig),
	}

}



func (manager *ClientHandler) send(message []byte) {
	for conn := range manager.clients {
		conn.send <- message

	}
}

func (manager *ClientHandler) showCharts(){
	for client := range manager.clients {
		client.send <- []byte("showCharts");
	}


}

func (manager *ClientHandler) start() {

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
