.PHONY: plugin
plugin:
	CGO_ENABLED=1 go build -buildmode=plugin -o ./build/missing_len.so ./plugin/missing_len.go
demo: plugin
	./build/`uname -s`_`uname -m`/golangci-lint run -Emissinglen ./testdata