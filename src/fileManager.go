package main

import (
	"os"
	"path/filepath"
	"github.com/nfnt/resize"
	"log"
	"bytes"
	"image/jpeg"
	"image"
)







const DirPath = "img/"
var FileList = make([]string,0)

func loadFiles () {

	filepath.Walk(DirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			FileList = append(FileList, info.Name())
		}
		return nil
	})

	for _, fileName := range FileList {
		println(fileName)
	}
}



func loadImgNoResize(manager *ClientManager) {

	buffer := new(bytes.Buffer)
	if currentPic + 1 >  len(FileList) {
		currentPic = 0
	}

	fimg, _ := os.Open(DirPath + FileList[currentPic])
	currentPic++
	imgReal, _, _ := image.Decode(fimg)
	if err := jpeg.Encode(buffer, imgReal, nil); err != nil {
		log.Println("unable to encode image.")
	}

	log.Println("foto en size real " , imgReal.Bounds().Max.X, imgReal.Bounds().Max.Y)

	manager.picBroadcast = make(chan map[*Client][]byte,len(manager.clients))

	go func() {
		for picClientMap := range manager.picBroadcast {
			for cli, pic := range picClientMap {
				cli.send <- pic

			}
		}

	}()

	manager.picture <- imgReal

}

func loadImg(manager *ClientManager) {

	buffer := new(bytes.Buffer)
	if currentPic + 1 >  len(FileList) {
		currentPic = 0
	}

	fimg, _ := os.Open(DirPath + FileList[currentPic])
	currentPic++
	imgReal, _, _ := image.Decode(fimg)
	if err := jpeg.Encode(buffer, imgReal, nil); err != nil {
		log.Println("unable to encode image.")
	}

	var resizedResolutionX uint
	var resizedResolutionY uint

	for client := range manager.clients{
		resizedResolutionX += uint (client.config.ResolutionWidth)
		if resizedResolutionY <  uint(client.config.Coordinate.Y) {
			resizedResolutionY = uint(client.config.Coordinate.Y + client.config.ResolutionHeight)
		}
	}


	m := resize.Resize(resizedResolutionX, resizedResolutionY , imgReal, resize.Lanczos3)

	log.Println("foto resaizada a " , resizedResolutionX, resizedResolutionY  )

	manager.picBroadcast = make(chan map[*Client][]byte,len(manager.clients))

	go func() {
		for picClientMap := range manager.picBroadcast {
			for cli, pic := range picClientMap {
				cli.send <- pic

			}
		}

	}()

	manager.picture <- m

}
