GO           ?= go
GOOS         ?= $(shell ${GO} env GOOS)
GOARCH       ?= $(shell ${GO} env GOARCH)
VERSION      ?= $(shell git describe --tags --abbrev=0 | sed -e 's/^v//')
BRANCH       ?= $(shell git rev-parse --abbrev-ref HEAD)
STATICCHECK  ?= $(shell $(GO) env GOPATH)/bin/staticcheck
REVIVE       ?= $(shell $(GO) env GOPATH)/bin/revive
GOFUMPT      ?= $(shell $(GO) env GOPATH)/bin/gofumpt
DOCKER       ?= docker
BINARY        = fastly-exporter
BINPKG        = ./cmd/fastly-exporter
SOURCE        = $(shell find . -name *.go)
DIST_DIR      = dist/v${VERSION}
DIST_BIN_FILE = ${BINARY}-${VERSION}.${GOOS}-${GOARCH}
DIST_ZIP_FILE = ${DIST_BIN_FILE}.tar.gz
DIST_BIN      = ${DIST_DIR}/${DIST_BIN_FILE}
DIST_ZIP      = ${DIST_DIR}/${DIST_ZIP_FILE}
DOCKER_TAG    = fastly-exporter:${VERSION}
DOCKER_ZIP    = ${DIST_DIR}/${BINARY}-${VERSION}.docker.tar.gz

${BINARY}: ${SOURCE} Makefile
	env CGO_ENABLED=0 ${GO} build -o ${BINARY} -ldflags="-X main.programVersion=${VERSION} -X github.com/prometheus/common/version.Version=${VERSION} -X github.com/prometheus/common/version.Branch=${BRANCH}" ${BINPKG}

${DIST_BIN}: ${DIST_DIR} ${SOURCE} Makefile
	env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} ${GO} build -o $@ -ldflags="-X main.programVersion=${VERSION} -X github.com/prometheus/common/version.Version=${VERSION} -X github.com/prometheus/common/version.Branch=${BRANCH}"  ${BINPKG}

${DIST_DIR}:
	mkdir -p $@

${DIST_ZIP}: ${DIST_BIN}
	tar -C ${DIST_DIR} -c -z -f ${DIST_ZIP} ${DIST_BIN_FILE}

.PHONY: docker-build
docker-build: ${SOURCE} Dockerfile
	${DOCKER} build --build-arg VERSION=${VERSION} --build-arg BRANCH=${BRANCH} --tag=${DOCKER_TAG} .

${DOCKER_ZIP}: docker-build
	${DOCKER} save --output=$@ ${DOCKER_TAG}

${STATICCHECK}:
	${GO} install honnef.co/go/tools/cmd/staticcheck@latest

${REVIVE}:
	${GO} install github.com/mgechev/revive@latest

${GOFUMPT}:
	${GO} install mvdan.cc/gofumpt@latest

.PHONY: lint
lint: ${STATICCHECK} ${REVIVE} ${GOFUMPT} ${SOURCE}
	${GO} vet ./...
	${STATICCHECK} ./...
	${REVIVE} ./...
	${GOFUMPT} -l -d -e .

.PHONY: test
test: ${SOURCE}
	${GO} test -race ./...

.PHONY: dist
dist: ${DIST_ZIP}

.PHONY: docker
docker: ${DOCKER_ZIP}

.PHONY: release
release: dist docker
