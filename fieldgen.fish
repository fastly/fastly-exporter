#!/usr/bin/env fish

cd cmd/fieldgen ; \
 and go run main.go > ../../pkg/gen/gen.go ; \
 and gofmt -w ../../pkg/gen/gen.go ; \
cd -
