NAME=ruleconverter
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CC ?= clang
upx ?= false

all: build


build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) go build -ldflags="-s -w" -o out/$(NAME)
	[ $(upx) ] && upx out/$(NAME)