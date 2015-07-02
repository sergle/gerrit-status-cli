ALL:
	mkdir -p bin src pkg
	GOPATH=`pwd` GOBIN="`pwd`/bin" go get

clean:
	rm -fr bin src pkg
