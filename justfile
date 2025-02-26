build:
	go build -ldflags="-s -w"

build-windows:
    go build -ldflags "-s -w -H=windowsgui"

test:
	go test -v ./...
