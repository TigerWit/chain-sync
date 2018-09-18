all: build

build: build-sync

build-sync:
	go build -o bin/sync ./cmd/sync
