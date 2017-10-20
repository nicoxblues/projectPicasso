package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	_ "golang.org/x/net/websocket"
	"gopkg.in/antage/eventsource.v1"

	_ "time"
	 _"encoding/json"
)



type deviceConfiguration struct{

	ResolutionHeight string  `json:"height"`
	ResolutionWidth string   `json:"width"`

	XCoordinate int  `json:"xValue"`
	YCoordinate int `json:"yValue"`

}

type Client struct {
	socket      *websocket.Conn
	clientID    string
	config      deviceConfiguration
	isConnected bool
	send   chan []byte
	prossPic   chan image.Image
	sendPic chan []byte
	graphicID int


}


type ClientManager struct {
	clients    map[*Client]bool
	picture    chan image.Image
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	picBroadcast chan []byte
}


var manager = ClientManager{
	picture:    make(chan image.Image),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
	picBroadcast: make(chan []byte),

}


var clientConn = make(map[*websocket.Conn]Client) //make(map[*websocket.Conn]bool) // connected clientConn




var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
				close(clieUnReg .send)
				delete(manager.clients, clieUnReg)
		//		jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
		//		manager.send(jsonMessage, cliReg)
			}
		case pic := <-manager.picture:
			for client := range manager.clients {
				select {
				case client.prossPic <- pic:
				default:
					close(client.send)
					delete(manager.clients, client)
				}
			}
			for clientToSend := range manager.clients{

			}


		}
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

	client := &Client{clientID:"test",

	config:deviceConfiguration{"1600","900",500,900},

	socket:conn,isConnected:true,graphicID:1}

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

			c.getByteCastPicture(picture)


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
		case _, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}



			//c.socket.WriteMessage(websocket.TextMessage, c.getByteCastPicture(picture))
		}
	}


	
}

func (c *Client) read() {
	
}

func (c *Client) getByteCastPicture(originalPic image.Image) []byte{


	return  []byte(getEncodeImage(c.getChunkImageForClient(originalPic)))

}

func (c *Client) getChunkImageForClient(image image.Image) (image.Image) {




	return nil
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

func writeImage(w http.ResponseWriter) {

	buffer := new(bytes.Buffer)
	fimg, _ := os.Open("img/imagen.jpg")
	imgReal, _, _ := image.Decode(fimg)
	if err := jpeg.Encode(buffer, imgReal, nil); err != nil {
		log.Println("unable to encode image.")
	}

	//w.Header().Set("Content-Type", "text/xml")
	w.Header().Set("Content-Type", "text/event-stream")
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

}

func main() {

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", wsHandler)

	es := eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/events", es)

	/*	go func() {
			for {
				e := string(getEncodeImage())
				es.SendEventMessage(e, "", "")
				log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
				time.Sleep(10 * time.Second)
			}
		}()

	*/

	err := http.ListenAndServe(":8080", nil) // set listen port
	if err == nil {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}

func chunkImage() {

	fimg, err := os.Open("img/imagen.jpg")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", "error to load file ", err)

	}

	defer fimg.Close()
	//  imgConfig, _, err := image.DecodeConfig(fimg)
	imgReal, _, err := image.Decode(fimg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", "error to decode file "+fimg.Name(), err)

	}

	bounds := imgReal.Bounds()
	chunkWidth := bounds.Max.X / 3 // determines the chunk width and height
	chunkHeight := bounds.Max.Y / 1
	fmt.Printf("%v\n%v\n", chunkHeight, chunkWidth)

	m0 := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
	draw.Draw(m0, image.Rect(0, 0, chunkWidth*2, chunkHeight), imgReal, image.Point{chunkWidth - (chunkWidth - 413), 0}, draw.Src)
	m1 := m0.SubImage(image.Rect(0, 0, 952, 1074)).(*image.RGBA)

	toimg, _ := os.Create("img/new2.jpg")
	defer toimg.Close()

	jpeg.Encode(toimg, m1, &jpeg.Options{jpeg.DefaultQuality})

	//  fmt.Printf("%v\n%v\n",chunkWidth,chunkHeight)
	// fmt.Printf("%v\n",m1)

}
