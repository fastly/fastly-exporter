# docker build --build-arg VERSION=$VERSION -t fastly-exporter:$VERSION .

FROM alpine:3.7

ARG VERSION

RUN apk add --no-cache ca-certificates

COPY dist/v${VERSION}/fastly-exporter-${VERSION}-linux-amd64 /fastly-exporter

ENTRYPOINT ["/fastly-exporter"]