// Package gen contains generated code that defines the rt.fastly.com API
// response, the Prometheus metrics we export, and the mapping between them.
package gen

//go:generate go run ../../cmd/fieldgen/main.go -metrics ../../cmd/fieldgen/exporter_metrics.json -fields ../../cmd/fieldgen/api_fields.json -mappings ../../cmd/fieldgen/mappings.json
