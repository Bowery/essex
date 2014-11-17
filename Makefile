DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps format
	@go get -d
	@go build
	./mercer

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

format:
	@echo "--> Running go fmt"
	@gofmt -w .

test: deps
	go test ./...

clean:
	rm -rf ~/.mercer
	rm mercer

.PHONY: all deps test format
