package main

import (
    "fmt"
    "net/http"
    //    mjpeg "github.com/marpie/go-mjpeg"
)

const URL = "http://user:user@10.10.10.50/mjpg/video.mjpg"

func main() {
    response, err := http.Get(URL)
    fmt.Println(err, response)
}
