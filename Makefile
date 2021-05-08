
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	go build -o node -ldflags "-X main.VERSION=$version -X 'main.BUILD_TIME=`date`' " main.go

.PHONY: node_0
node_0:
	go build -o node_0 -ldflags "-X main.VERSION=$version -X 'main.BUILD_TIME=`date`' " main.go

.PHONY: node_1
node_1:
	go build -o node_1 -ldflags "-X main.VERSION=$version -X 'main.BUILD_TIME=`date`' " main.go