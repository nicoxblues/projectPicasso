package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"


)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	socket    *websocket.Conn
	clientID   string
	config    deviceConfiguration
	send      chan []byte
	prossPic  chan image.Image
	sendPic   chan []byte
	clientChart *chartConfig
	manager   *ClientHandler
}



func wsHandler(manager *ClientHandler, w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "oh no ! :(", 403)
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



	client := &Client{
		config: deviceConfiguration{height, width, image.Point{coorX, coorY}},
		socket: conn,  prossPic: make(chan image.Image), send: make(chan []byte), manager: manager}

	client.loadClientConfig()

	client.manager.charConf[client.clientChart.ChartID] = client.clientChart


	client.manager.register <- client
	//conn.WriteJSON(client.clientChart)
	conn.WriteMessage(websocket.TextMessage,[]byte(client.clientChart.HtmlDivRoot))

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

func (c *Client) loadClientConfig() {
	if c.clientChart == nil {
		c.clientChart  = c.manager.serverConf.nextChart()
		c.clientID= c.clientChart.ChartID
	}


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


