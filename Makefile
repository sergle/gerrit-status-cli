ALL:
	mkdir -p src pkg
	GOPATH=`pwd` GOBIN="`pwd`" GO_ENABLED=0 GOOS=linux go get -a -tags netgo -ldflags '-w'

clean:
	rm -fr src pkg

static:
	GOPATH=`pwd` GO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
mac:
	GOPATH=`pwd` GO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -ldflags '-w' .
