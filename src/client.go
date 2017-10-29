package main

import (

	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	"fmt"
	"image/draw"
	"bytes"
	"encoding/base64"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type Client struct {
	socket      *websocket.Conn
	clientID    string
	config      deviceConfiguration
	isConnected bool
	send        chan []byte
	prossPic    chan image.Image
	sendPic     chan []byte
	graphicID   int
	manager *ClientManager
}




func wsHandler(manager *ClientManager , w http.ResponseWriter  , r *http.Request) {

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	queryValues := r.URL.Query()

	width, _ := strconv.Atoi(queryValues.Get("width"))
	height, _ := strconv.Atoi(queryValues.Get("height"))
	coorX, _ := strconv.Atoi(queryValues.Get("coordenadasX"))
	coorY, _ := strconv.Atoi(queryValues.Get("coordenadasY"))

	client := &Client{clientID: "test",
		config: deviceConfiguration{height, width, image.Point{coorX, coorY}},
		socket: conn, isConnected: true, graphicID: 1, prossPic: make(chan image.Image), send: make(chan []byte),manager:manager}

	//manager.register <- client

	client.manager.register <- client

	go client.processPic()
	go client.read()
	go client.write()



}

func (c *Client) processPic() {

	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case picture, ok := <-c.prossPic:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			var clientPicMap = make(map[*Client][]byte)
			clientPicMap[c] = c.getByteCastPicture(&picture)
			c.manager.picBroadcast <- clientPicMap

		}
	}

}

func (c *Client) write() {

	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case mesagge, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, mesagge)
		}
	}

}

func (c *Client) read() {

	defer func() {
		c.manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, _, err := c.socket.ReadMessage() // por ahora solo me importa saber si el cliente pudo responder
		if err != nil {
			c.manager.unregister <- c
			c.socket.Close()
			break
		}





	}

}

func (c *Client) getByteCastPicture(originalPic *image.Image) []byte {

	return []byte(getEncodeImage(c.getChunkImageForClient(originalPic)))

}

func (c *Client) getChunkImageForClient(originImage *image.Image) image.Image {

	chunkWidth := c.config.ResolutionWidth
	chunkHeight := c.config.ResolutionHeight
	fmt.Printf("%v cantidad de clientes conectados \n", len(c.manager.clients))

	m0 := image.NewRGBA(image.Rect(chunkWidth, chunkHeight, 0, 0))
	//m1 := image.NewRGBA((*originImage).Bounds())

	draw.Draw(m0, (*originImage).Bounds(), *originImage, c.config.Coordinate, draw.Src)
	//draw.Draw(m1, (*originImage).Bounds(), *originImage, image.Point{0,0}, draw.Src)

	return m0 //resize.Resize(uint(chunkWidth), uint(chunkHeight),m0,resize.Lanczos3)
}


func getEncodeImage(image image.Image) string {

	buffer := new(bytes.Buffer)

	if err := jpeg.Encode(buffer, image, nil); err != nil {
		log.Println("unable to encode image.")
	}

	encodedString := base64.StdEncoding.EncodeToString(buffer.Bytes())
	//w.Header().Set("Content-Type", "text/xml")

	return encodedString

}

