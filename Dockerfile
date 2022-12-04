ARG GO_VERSION=1.19
ARG XX_VERSION=1.1.2

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine as builder

# Copy the build utilities.
COPY --from=xx / /

ARG TARGETPLATFORM

WORKDIR /workspace

# copy api submodule
COPY api/ api/

# copy modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

ENV GOPROXY=https://goproxy.io,direct

# cache modules
RUN go mod download

# copy source code
COPY main.go main.go
COPY controllers/ controllers/
COPY internal/ internal/

# build
ENV CGO_ENABLED=0
RUN xx-go build -trimpath -a -o notification-controller main.go

FROM alpine:3.16

LABEL org.opencontainers.image.source="https://github.com/fluxcd/notification-controller"

RUN apk add --no-cache ca-certificates tini

COPY --from=builder /workspace/notification-controller /usr/local/bin/

USER 65534:65534

ENTRYPOINT [ "/sbin/tini", "--", "notification-controller" ]
