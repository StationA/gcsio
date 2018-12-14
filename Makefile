tools:
	@go install github.com/golang/dep/cmd/dep

deps: tools
	dep ensure

build: tools
	@go build -o target/gcsio ./...

install: build
	@go install ./...

target:
	mkdir -p target

release: build target
	@CGO_ENABLED=0 go build -a -o target/gcsio ./cmd/gcsio

release-all: build target
	@CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -o target/gcsio.darwin-386 ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o target/gcsio.darwin-amd64 ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -o target/gcsio.linux-386 ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o target/gcsio.linux-amd64 ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -o target/gcsio.linux-arm ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o target/gcsio.linux-arm64 ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -o target/gcsio.windows-386.exe ./cmd/gcsio
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o target/gcsio.windows-amd64.exe ./cmd/gcsio

clean:
	@rm -rf target

.PHONY: tools build install release release-all clean
