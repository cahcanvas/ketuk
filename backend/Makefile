.PHONY: run tidy test bundle-api

run:
	go run ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

bundle-api:
	npx --yes @redocly/cli bundle api/openapi.yaml -o api/openapi.bundled.yaml
