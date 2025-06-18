FROM golang:latest AS builder
ARG VERSION=0.0.0-unknown
ARG BRANCH=main

RUN groupadd -r fastly-exporter
RUN useradd -r -g fastly-exporter fastly-exporter

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

ADD .git .git
ADD cmd cmd
ADD pkg pkg

RUN env CGO_ENABLED=0 go build \
	-a \
	-ldflags="-X main.programVersion=$VERSION -X github.com/prometheus/common/version.Version=$VERSION -X github.com/prometheus/common/version.Branch=$BRANCH" \
	-o /fastly-exporter \
	./cmd/fastly-exporter

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /fastly-exporter /fastly-exporter

USER fastly-exporter

EXPOSE 8080

ENTRYPOINT ["/fastly-exporter", "-listen=0.0.0.0:8080"]

LABEL org.opencontainers.image.source="https://github.com/fastly/fastly-exporter/"
LABEL org.opencontainers.image.description="Fastly Prometheus Exporter container image"
LABEL org.opencontainers.image.licenses="Apache-2.0"
