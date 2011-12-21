package main

import (
    "fmt"
    mjpeg "github.com/marpie/go-mjpeg"
    "image"
    "io"
    "net/http"
    "time"
)

const URL = "http://user:user@10.10.10.50/mjpg/video.mjpg"

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
        fmt.Println("New Image:", img.ColorModel())
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
