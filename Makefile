.PHONY: help
help:
	@echo 'Makefile for `godocjson` project'
	@echo ''
	@echo 'Usage:'
	@echo '   make godocjson-darwin-amd64    build binary'
	@echo '   make golembic.json             re-build JSON for docs'
	@echo ''

godocjson-darwin-amd64: main.go util.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./godocjson-darwin-amd64 .

golembic.json: godocjson-darwin-amd64
	./godocjson-darwin-amd64 "${GOPATH}/src/github.com/dhermes/golembic" > golembic.json
