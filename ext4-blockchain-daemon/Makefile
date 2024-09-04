.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o ./bin/ext4-chain-daemon ./cmd/ext4-chain-daemon

