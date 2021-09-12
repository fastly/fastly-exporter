#!/usr/bin/env fish

pushd cmd/fieldgen

go run main.go > ../../pkg/gen/gen.go
and gofmt -w ../../pkg/gen/gen.go

popd

