package mjpeg

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "runtime/pprof"
    "testing"
)

const (
    testfile        = "./testdata/video.mjpg"
    useProfiling    = true
    profilingCpu    = "./profile_cpu.prof"
    profilingMemory = "./profile_mem.prof"
)

var sampleData *bytes.Buffer

func setup() bool {
    inFile, err := os.Open(testfile)
    defer inFile.Close()
    if err != nil {
        return false
    }
    slice := make([]byte, 1914928)
    if _, rErr := inFile.Read(slice); rErr != nil {
        return false
    }
    sampleData = bytes.NewBuffer(slice)
    return sampleData != nil
}

func startProfiling() {
    if !useProfiling {
        return
    }

    f, _ := os.Create(profilingCpu)
    pprof.StartCPUProfile(f)
}

func stopProfiling() {
    if !useProfiling {
        return
    }
    pprof.StopCPUProfile()
}

func TestReadHeader(t *testing.T) {
    if !setup() {
        t.Error("setup() failed")
    }

    h, err := readHeader(bufio.NewReader(sampleData))
    switch {
    case err != nil:
        t.Errorf("readHeader error: %v", err)
    case h.boundary != "myboundary":
        t.Errorf("boundary mismatch: '%v'", h.boundary)
    case h.motion_event != 1:
        t.Errorf("motion_event mismatch: '%v'", h.motion_event)
    case h.content_type != "image/jpeg":
        t.Errorf("content_type mismatch: '%v'", h.content_type)
    case h.content_length != 14724:
        t.Errorf("content_length mismatch: '%v'", h.content_length)
    }
}

func TestDecode(t *testing.T) {
    startProfiling()
    defer stopProfiling()
    if !setup() {
        t.Error("setup() failed.")
    }

    i := 0
    for {
        img, err := Decode(sampleData)
        if err != nil && err != io.EOF {
            t.Errorf("Decode failed with error: %v (i == %d)", err, i)
        }
        if img == nil {
            break
        }
        i++
    }
}
