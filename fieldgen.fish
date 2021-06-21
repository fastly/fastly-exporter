#!/usr/bin/env fish

cd cmd/fieldgen ; \
go run main.go > ../../pkg/gen/gen.go
gofmt -w ../../pkg/gen/gen.go
cd -
