GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

.PHONY: build
dms-server:
	CGO_ENABLED=0 go build -o ./bin/dms-server ./cmd/dms/dms-server.go

.PHONY: docker-build
docker-build: build
docker-build:
	@docker build -t dms-server:latest ./ 
