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
	"os"

	_ "encoding/json"
	"strconv"
)

type deviceConfiguration struct {
	ResolutionHeight int
	ResolutionWidth  int

	Coordinate image.Point
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
}

type ClientManager struct {
	clients      map[*Client]bool
	picture      chan image.Image
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	picBroadcast chan map[*Client][]byte
}

var manager = ClientManager{
	picture:      make(chan image.Image),
	broadcast:    make(chan []byte),
	register:     make(chan *Client),
	unregister:   make(chan *Client),
	clients:      make(map[*Client]bool),

}

//var clientConn = make(map[*websocket.Conn]Client) //make(map[*websocket.Conn]bool) // connected clientConn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

				/*default:
				close(client.send)
				delete(manager.clients, client)*/
			}


		/*	case   messageMap := <-manager.picBroadcast:
			for cli := range manager.clients {
				select {
					case cli.send <- messageMap[cli]:

					default:
						close(cli.send)
						delete(manager.clients, cli)

				}
			}*/
		}




			/*for cli := range manager.clients {
				select {
				case cli.send <- PicClientMap[cli]:

				default:
					close(cli.send)
					delete(manager.clients, cli)

				}
			}*/
			//for clientToSend := range manager.clients{

		//}

	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) //websocket.Upgrade(w, r, w.Header(), 1024, 1024)

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
		socket: conn, isConnected: true, graphicID: 1, prossPic: make(chan image.Image), send: make(chan []byte)}

	manager.register <- client

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
			manager.picBroadcast <- clientPicMap
			//c.send <- c.getByteCastPicture(&picture)
			//c.socket.WriteMessage(websocket.TextMessage, c.getByteCastPicture(picture))
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

}

func (c *Client) getByteCastPicture(originalPic *image.Image) []byte {

	return []byte(getEncodeImage(c.getChunkImageForClient(originalPic)))

}

func (c *Client) getChunkImageForClient(originImage *image.Image) image.Image {

	//fimg, err := os.Open("img/imagen.jpg")



	//bounds := imgReal.Bounds()
	chunkWidth := c.config.ResolutionWidth
	chunkHeight := c.config.ResolutionHeight
	fmt.Printf("%v\n%v\n cantidad de clientes conectados", len(manager.clients)  )

	//rec := imgReal.Bounds()
	m0 := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
	//draw.Draw(m0, image.Rect(0, 0, rec.Max.X, rec.Max.Y), imgReal, c.config.Coordinate, draw.Src)
	draw.Draw(m0, image.Rect(0, 0, 7680, 4800), *originImage, c.config.Coordinate, draw.Src)

	//m1 := m0.SubImage(image.Rect(0, 0, chunkWidth, chunkHeight)).(*image.RGBA)

	return m0
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

func loadImg(w http.ResponseWriter, r *http.Request) {

	buffer := new(bytes.Buffer)
	fimg, _ := os.Open("img/imagen.jpg")
	imgReal, _, _ := image.Decode(fimg)
	if err := jpeg.Encode(buffer, imgReal, nil); err != nil {
		log.Println("unable to encode image.")
	}

	manager.picBroadcast = make(chan map[*Client][]byte,len(manager.clients))

	go func() {
		for {
			for picClientMap := range manager.picBroadcast {
				for cli, pic := range picClientMap {
					cli.send <- pic

				}
			}
		}
	}()

	manager.picture <- imgReal

}

func main() {

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/loadImage", loadImg)

	go manager.start()

	err := http.ListenAndServe(":8080", nil) // set listen port
	if err == nil {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
