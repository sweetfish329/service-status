.PHONY: build
build:
	mkdir -p build
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.33.0
	GOOS=js GOARCH=wasm go build -o build/app.wasm main.go html.go
