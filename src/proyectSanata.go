package main

import (
	"image"
	"log"
	"net/http"

)

type deviceConfiguration struct {
	ResolutionHeight int
	ResolutionWidth  int

	Coordinate image.Point
}

func main() {

	imgWrapper := initManager()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	clientHandler := newClientHandler()

	go clientHandler.start()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsHandler(clientHandler, writer, request)
	})

	http.HandleFunc("/loadImage", func(writer http.ResponseWriter, request *http.Request) {
		imgWrapper.shouldResizeImage = true
		imgWrapper.loadImg(clientHandler)
	})
	http.HandleFunc("/loadImage2", func(writer http.ResponseWriter, request *http.Request) {
		imgWrapper.shouldResizeImage = false
		imgWrapper.loadImg(clientHandler)
	})
	http.HandleFunc("/showChart", func(writer http.ResponseWriter, request *http.Request) {
		clientHandler.showCharts()
	})
	http.HandleFunc("/resetFinneg", func(writer http.ResponseWriter, request *http.Request) {
		clientHandler.resetFin()
	})

	http.HandleFunc("/manager", func(w http.ResponseWriter, request *http.Request) {
			mainPageLoader(w,clientHandler)
	})






	err := http.ListenAndServe(":8080", nil) // set listen port
	if err == nil {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
