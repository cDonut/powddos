default: build-server-image build-client-image start

.PHONY: build-server-image
build-server-image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o powserver ./cmd/server/main.go
	docker build -t powserver . -f ./docker/server/Dockerfile
	rm powserver

.PHONY: build-client-image
build-client-image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o powclient ./cmd/client/main.go
	docker build -t powclient . -f ./docker/client/Dockerfile
	rm powclient

.PHONY: start
start:
	docker-compose -f docker/docker-compose.yml up
