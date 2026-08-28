.PHONY: build test lint tidy docker

build:
	cd src && go build -o ../bin/kube-scheduler ./cmd/scheduler

test:
	cd src && go test ./... -race -cover

lint:
	cd src && golangci-lint run

tidy:
	cd src && go mod tidy

docker:
	docker buildx build --platform linux/amd64,linux/arm64 -t kube-scheduler:dev .
