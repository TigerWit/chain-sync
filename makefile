all: build

build: build-explorer build-sync build-incr-sync

build-explorer:
	go build -o bin/explorer .

build-sync:
	go build -o bin/sync ./cmd/sync

build-incr-sync:
	go build -o bin/incr-sync ./cmd/incr-sync