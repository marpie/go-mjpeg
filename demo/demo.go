// Copyright 2011 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This demo program shows a basic example of the MJPEG library.
package main

import (
	"fmt"
	mjpeg "github.com/marpie/go-mjpeg"
	"image"
	"io"
	"net/http"
	"time"
)

const URL = "http://user:user@localhost:5050/mjpg/video.mjpg"

// processHttp receives the HTTP data and tries to decodes images. The images 
// are sent through a chan for further processing.
func processHttp(response *http.Response, nextImg chan *image.Image, quit chan bool) {
	defer response.Body.Close()
	for {
		select {
		case <-quit:
			close(nextImg)
			return
		default:
			img, err := mjpeg.Decode(response.Body)
			if err == io.EOF {
				close(nextImg)
				return
			}
			if err != nil {
				fmt.Println(err)
			}
			if img != nil {
				nextImg <- img
			}
		}
	}
}

// processImage receives images through a chan and prints the dimensions.
func processImage(nextImg chan *image.Image, quit chan bool) {
	for i := 0; i < 10; i++ {
		i, ok := <-nextImg
		if !ok {
			break
		}
		img := *i
		if *i == nil {
			continue
		}
		fmt.Println("New Image:", img.Bounds())
	}
	quit <- true
}

func main() {
	response, err := http.Get(URL)
	if err != nil {
		return
	}
	nextImg := make(chan *image.Image, 30)
	quit := make(chan bool)
	fmt.Println("Waiting for images to process...")
	go processImage(nextImg, quit)
	go processHttp(response, nextImg, quit)
	time.Sleep(20 * time.Second)
}
