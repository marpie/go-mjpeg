// Copyright 2011 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mjpeg implements the Motion JPEG format of the EDIMAX webcam.
// It may be compatible with other MJPEG implementations.
package mjpeg

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"strconv"
	"strings"
)

// header contains the fields that preceed a given image.
type header struct {
	boundary       string
	motion_event   int
	content_type   string
	content_length int
}

// readString tries to read a string until delim is found. The delim byte 
// won't be returned.
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
	return "", errors.New("Unknown error")
}

// readHeader is used to find and return a correct MJPEG header.
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
			return nil, fmt.Errorf("Not enough data available (2 needed - got %v: %X).", n, data[0])
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
			return nil, errors.New("Not a valid key/value pair.")
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

// Decode returns the next Image found in the MJPEG stream.
func Decode(inReader io.Reader) (img *image.Image, out_err error) {
	// read header
	h, out_err := readHeader(inReader)
	if out_err != nil {
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

	return nil, errors.New("Unknown error")
}
