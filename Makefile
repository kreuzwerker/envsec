VERSION := $(shell git describe --abbrev=7 --always)
REPO := envsec
USER := kreuzwerker
FLAGS := "-X=main.build=$(VERSION) -X=main.version=$(VERSION)"

.PHONY: build clean release retract

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags $(FLAGS) -o build/linux-amd64/ep bin/es.go
	GOOS=linux GOARCH=arm go build -ldflags $(FLAGS) -o build/linux-arm/ep bin/es.go
	GOOS=darwin GOARCH=amd64 go build -ldflags $(FLAGS) -o build/darwin-amd64/ep bin/es.go

clean:
	rm -rf build