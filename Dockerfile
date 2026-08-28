# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/cmd/ cmd/
COPY src/pkg/ pkg/

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/kube-scheduler ./cmd/scheduler

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/kube-scheduler /kube-scheduler
USER 65532:65532
ENTRYPOINT ["/kube-scheduler"]
