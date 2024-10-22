###############################################
FROM --platform=$BUILDPLATFORM golang:1.23@sha256:a7f2fc9834049c1f5df787690026a53738e55fc097cd8a4a93faa3e06c67ee32 AS builder

WORKDIR /app

COPY . .

ARG TARGETOS TARGETARCH
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

RUN go build -o bin/p2p cmd/p2p/main.go

###############################################
FROM busybox AS busybox
RUN addgroup -S p2pgroup && adduser -S -G p2pgroup p2puser

###############################################
FROM debian:stable-slim AS certs
RUN apt-get update && apt-get install -y ca-certificates

###############################################
FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/p2p .

COPY --from=busybox /etc/passwd /etc/passwd
COPY --from=busybox /etc/group /etc/group

COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=certs /etc/ssl/private /etc/ssl/private
COPY --from=certs /usr/share/ca-certificates /usr/share/ca-certificates

USER p2puser:p2pgroup

ENTRYPOINT ["/app/p2p"]
CMD []