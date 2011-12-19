package mjpeg

import (
    "bufio"
    "image"
    "image/jpeg"
    "io"
    "strconv"
    "strings"
)

type header struct {
    boundary       string
    motion_event   int
    content_type   string
    content_length int
}

func readHeader(inReader io.Reader) (h *header, ok bool) {
    var err error
    h = new(header)
    r := bufio.NewReader(inReader)

    h.boundary, err = r.ReadString('\n')
    if err != nil || h.boundary[0] != '-' || h.boundary[1] != '-' {
        return nil, false
    }
    h.boundary = strings.TrimSpace(h.boundary)

    for {
        kv_str, err := r.ReadString('\n')
        if err != nil || len(kv_str) == 0 {
            break
        }

        kv := strings.SplitN(kv_str, ":", 2)
        if len(kv) != 2 {
            break
        }
        value := strings.TrimSpace(kv[1])
        switch kv[0] {
        case "Motion-Event":
            h.motion_event, _ = strconv.Atoi(value)
        case "Content-Type":
            h.content_type = value
        case "Content-Length":
            h.content_length, _ = strconv.Atoi(value)
        }
    }
    return h, true
}

func Decode(r io.Reader) (img *image.Image, out_ok bool) {
    // read header
    h, ok := readHeader(r)
    if !ok || h.content_length < 1 {
      return nil, false
    }

    // read image content
    switch h.content_type {
      case "image/jpeg":
        img, err 
      default:
        return nil, false
    }

    return nil, false
}
