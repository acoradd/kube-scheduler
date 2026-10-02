# syntax=docker/dockerfile:1
FROM golang:1.27 AS build
WORKDIR /build

COPY src/go.mod src/go.sum src/
COPY src/cmd/ src/cmd/
COPY src/pkg/ src/pkg/
COPY hack/ hack/

# Kubernetes release to build against (e.g. 1.35.9); defaults to src/go.mod.
ARG K8S_VERSION
RUN --mount=type=cache,target=/go/pkg/mod \
    if [ -n "${K8S_VERSION}" ]; then bash hack/set-k8s-version.sh "${K8S_VERSION}"; fi && \
    cd src && go mod download

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd src && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/kube-scheduler ./cmd/scheduler

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/kube-scheduler /kube-scheduler
USER 65532:65532
ENTRYPOINT ["/kube-scheduler"]
