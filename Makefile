default:
	# build
	go mod tidy
	command -v goimports || go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	go vet ./...
	go run golang.org/x/tools/cmd/deadcode@latest ./...
	go test -cover ./...
	go install .
	# docs
	rm -Rf out docs/reference
	go run ./hack update-schema
	command -v jsonschema2md || npm install -g @adobe/jsonschema2md
	jsonschema2md -d schema -f yaml -o docs/reference

.PHONY: release
release:
	./release
