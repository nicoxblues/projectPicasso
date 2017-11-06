package main

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

type fileHandler struct {
	fileNameCacheList *[]string
	currentFile       int
}

type imageWrapper struct {
	shouldResizeImage bool
	img               image.Image
	fileList          fileHandler
}

const DirPath = "img/"

func (fh *fileHandler) getImage() (image.Image, error) {

	cacheFile := *fh.fileNameCacheList

	buffer := new(bytes.Buffer)
	if fh.currentFile+1 > len(cacheFile) {
		fh.currentFile = 0
	}

	fimg, _ := os.Open(DirPath + cacheFile[fh.currentFile])
	fh.currentFile++
	imgReal, _, _ := image.Decode(fimg)
	if err := jpeg.Encode(buffer, imgReal, nil); err != nil {
		log.Println("unable to encode image.")
		return nil, err
	}

	return imgReal, nil

}

func initManager() *imageWrapper {

	var fileCache []string
	wrapper := &imageWrapper{fileList: fileHandler{&fileCache, 0}}

	filepath.Walk(DirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileCache = append(fileCache, info.Name())
		}
		return nil
	})

	for _, fileName := range fileCache {
		println(fileName)
	}

	return wrapper

}

func (wrapper *imageWrapper) loadImg(manager *ClientHandler) {

	imgReal, err := wrapper.fileList.getImage()
	log.Println("foto en size real ", imgReal.Bounds().Max.X, imgReal.Bounds().Max.Y)
	if err == nil {

	}

	var resizeResolutionX uint
	var resizeResolutionY uint

	if wrapper.shouldResizeImage {

		for client := range manager.clients {
			resizeResolutionX += uint(client.config.ResolutionWidth)
			if resizeResolutionY < uint(client.config.Coordinate.Y) {
				resizeResolutionY = uint(client.config.Coordinate.Y + client.config.ResolutionHeight)
			}
			if resizeResolutionY == 0{
				resizeResolutionY = uint(client.config.ResolutionHeight);
			}
		}

		imgReal = resize.Resize(resizeResolutionX, resizeResolutionY, imgReal, resize.Lanczos3)

		log.Println("foto resaizada a ", resizeResolutionX, resizeResolutionY)
	}

	manager.picBroadcast = make(chan map[*Client][]byte, len(manager.clients))

	go func() {
		for picClientMap := range manager.picBroadcast {
			for cli, pic := range picClientMap {
				cli.send <- pic

			}
		}

	}()

	manager.picture <- imgReal

}
