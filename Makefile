.PHONY: plugin
plugin:
	CGO_ENABLED=1 go build -buildmode=plugin -o ./build/`uname -s`_`uname -m`/missing_len.so ./plugin/missing_len.go
	cp ./build/`uname -s`_`uname -m`/missing_len.so ./build/missing_len.so

.PHONY: ci-lint-demo
ci-lint-demo: plugin
	./build/`uname -s`_`uname -m`/golangci-lint run -Emissinglen ./testdata

.PHONY: lint
lint: plugin
	./build/`uname -s`_`uname -m`/golangci-lint run -Emissinglen ./pkg

.PHONY: build
build: plugin
	go build -a -o ./build/`uname -s`_`uname -m`/linter ./cmd/main.go

.PHONY: demo
demo: build
	./build/`uname -s`_`uname -m`/linter --file ./testdata/sample.go

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: init
init: tidy

.PHONY: all
all: tidy lint build