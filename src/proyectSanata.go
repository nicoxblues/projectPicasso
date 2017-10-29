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



var currentPic = 0

func main() {

	loadFiles()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	manager := newManager()

	go manager.start()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsHandler(manager,writer,request)
	})

	http.HandleFunc("/loadImage", func(writer http.ResponseWriter, request *http.Request) {
		loadImg(manager)
	})
	http.HandleFunc("/loadImage2", func(writer http.ResponseWriter, request *http.Request) {
		loadImgNoResize(manager)
	})



	err := http.ListenAndServe(":8080", nil) // set listen port
	if err == nil {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
