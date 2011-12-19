package mjpeg

import (
    "strings"
    "testing"
)

const SampleHeader = `--myboundary
Motion-Event: 1
Content-Type: image/jpeg
Content-Length: 14724

`

var fakeData = strings.NewReader(SampleHeader)

func TestReadHeader(t *testing.T) {
    h, ok := readHeader(fakeData)
    switch {
    case !ok:
        t.Errorf("readHeader not ok!")
    case h.boundary != "--myboundary":
        t.Errorf("boundary mismatch: '%v'", h.boundary)
    case h.motion_event != 1:
        t.Errorf("motion_event mismatch: '%v'", h.motion_event)
    case h.content_type != "image/jpeg":
        t.Errorf("content_type mismatch: '%v'", h.content_type)
    case h.content_length != 14724:
        t.Errorf("content_length mismatch: '%v'", h.content_length)
    }
}
