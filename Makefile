ALL:
	mkdir -p src pkg
	GOPATH=`pwd` GOBIN="`pwd`" go get

clean:
	rm -fr src pkg

static:
	GOPATH=`pwd` GO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
