COMMIT := $(shell git describe --dirty --always)
LDFLAGS := "-s -w -X main.GitCommit=$(COMMIT)"
DOCKER_IMAGE ?= networkop/xdp-xconnect

generate:
	go generate ./...

build:
	go build -o xdp-xconnect main.go 
	
lint:
	golangci-lint run

test:
	go test -race ./...  -v

docker: Dockerfile test
	docker buildx build --push \
	--platform linux/amd64 \
	--build-arg LDFLAGS=$(LDFLAGS) \
	-t $(DOCKER_IMAGE):$(COMMIT) \
	-t $(DOCKER_IMAGE):latest .