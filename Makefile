

.PHONY: install tidy build dist fmt

build: tidy fmt
	@go build ./...

install: tidy
	@go install ./...

dist:
	@goreleaser --skip-publish --snapshot --rm-dist

tidy:
	@go get
	@go mod tidy

fmt:
	@gofumpt -w .