package mjpeg

import (
    "bytes"
    "fmt"
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

type ErrorString string

func (e ErrorString) Error() string { return string(e) }
func NewError(s string) error       { return ErrorString(s) }

func readString(inReader io.Reader, delim byte) (string, error) {
    var b = make([]byte, 1)
    buffer := bytes.NewBuffer(nil)
    for {
        n, err := inReader.Read(b)
        if err != nil || n < 1 {
            return "", err
        }
        if b[0] == delim {
            return strings.TrimSpace(buffer.String()), nil
        } else {
            buffer.Write(b)
        }
    }
    return "", NewError("Unknown error")
}

func readHeader(inReader io.Reader) (h *header, out_err error) {
    // search for boundary
    data := make([]byte, 2)
    for {
        n, err := inReader.Read(data)
        switch {
        case err != nil:
            return nil, err
        case n == 1 && ((data[0] == 0x0A) || (data[0] == 0x0D)):
            continue
        case n < 2:
            return nil, NewError(fmt.Sprintf("Not enough data available (2 needed - got %v: %X).", n, data[0]))
        }

        if data[0] == '-' && data[1] == '-' {
            break
        }
    }

    // populate header
    h = new(header)
    h.boundary, out_err = readString(inReader, '\n')
    if out_err != nil {
        return nil, out_err
    }

    for {
        line, err := readString(inReader, '\n')
        if err != nil {
            return nil, err
        }

        if line == "" {
            break
        }

        kv := strings.Split(line, ": ")
        if len(kv) != 2 {
            return nil, NewError("Not a valid key/value pair.")
        }

        switch kv[0] {
        case "Motion-Event":
            h.motion_event, out_err = strconv.Atoi(kv[1])
            if out_err != nil {
                return nil, out_err
            }
        case "Content-Type":
            h.content_type = kv[1]
        case "Content-Length":
            h.content_length, out_err = strconv.Atoi(kv[1])
            if out_err != nil {
                return nil, out_err
            }
        }

    }
    return h, nil
}

func Decode(inReader io.Reader) (img *image.Image, out_err error) {
    // read header
    h, out_err := readHeader(inReader)
    if out_err != nil { //|| h.content_length < 1 {
        return nil, out_err
    }

    switch h.content_type {
    case "image/jpeg":
        jpg, out_err := jpeg.Decode(inReader)
        if out_err != nil {
            return nil, out_err
        }
        return &jpg, nil
    }

    return nil, NewError("Unknown error")
}
