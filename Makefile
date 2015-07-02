ALL:
	mkdir -p src pkg
	GOPATH=`pwd` GOBIN="`pwd`" go get

clean:
	rm -fr src pkg
