# Based off https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

# Accept the Go version for the image to be set as a build argument
ARG GO_VERSION=1.12.7

# First stage: build the executable
FROM golang:${GO_VERSION}-alpine AS builder

# Fetch go modules
# https://github.com/golang/go/wiki/Modules
# https://dev.to/maelvls/why-is-go111module-everywhere-and-everything-about-go-modules-24k
ENV GO111MODULE=on

# ca-certificates for calls to HTTPS endpoints
# git for fetching the dependencies
RUN apk add --no-cache \
	ca-certificates \
	git

# Create appuser
RUN adduser -D -g '' appuser

# Get code
COPY . $GOPATH/src/github.com/peterbourgon/fastly-exporter/
WORKDIR $GOPATH/src/github.com/peterbourgon/fastly-exporter/

# Build the binary
RUN CGO_ENABLED=0 go build \
	-a \
	-ldflags="-X main.programVersion=$(git describe | sed -e 's/^v//')" \
	-o /go/bin/fastly-exporter \
	./cmd/fastly-exporter

# Second stage: build the container
FROM scratch

# Copy dependencies
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /go/bin/fastly-exporter /go/bin/fastly-exporter

USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/fastly-exporter", "-endpoint", "http://0.0.0.0:8080/metrics"]

