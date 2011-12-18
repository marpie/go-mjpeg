include $(GOROOT)/src/Make.inc

TARG=github.com/marpie/go-mjpeg
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	mjpeg.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w mjpeg.go
	${GOFMT} -w mjpeg_test.go
