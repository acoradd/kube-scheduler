.PHONY: build test lint tidy docker k8s-version

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

# Retarget src/go.mod onto another Kubernetes release, e.g. make k8s-version K8S_VERSION=1.35.9
k8s-version:
	bash hack/set-k8s-version.sh $(K8S_VERSION)
