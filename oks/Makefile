all: help

.PHONY: help
help:
	@echo "help:"
	@echo "- make gen  : regenerate SDK"
	@echo "- make test : run all tests"

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: generate
generate:
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.6.0
	go get k8s.io/code-generator/cmd/deepcopy-gen@v0.36.2
	go generate ./...

.PHONY: gen
gen: generate tidy

.PHONY: test
test: gen
	go test -v ./examples
