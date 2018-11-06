FROM golang:1.10-alpine3.7

ENV APP_DIR=/opt/app

RUN mkdir -p ${APP_DIR}

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR ${APP_DIR}

RUN go get -d -v github.com/peterbourgon/fastly-exporter
RUN go install -v github.com/peterbourgon/fastly-exporter

ADD entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
