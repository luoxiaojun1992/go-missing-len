.PHONY: plugin
plugin:
	CGO_ENABLED=1 go build -buildmode=plugin -o ./build/`uname -s`_`uname -m`/missing_len.so ./plugin/missing_len.go
	cp ./build/`uname -s`_`uname -m`/missing_len.so ./build/missing_len.so

.PHONY: demo
demo: plugin
	./build/`uname -s`_`uname -m`/golangci-lint run -Emissinglen ./testdata

.PHONY: lint
lint: plugin
	./build/`uname -s`_`uname -m`/golangci-lint run -Emissinglen ./pkg

.PHONY: build
build:
	go build -a -o ./build/`uname -s`_`uname -m`/linter ./cmd/main.go