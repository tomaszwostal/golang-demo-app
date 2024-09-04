ENV	:= $(PWD)/.env
include $(ENV)
OS := $(shell uname -s)

# Change these variables as necessary.
main_package_path = ./cmd/demo
binary_name = demo

.PHONY: audit
audit: test 
	go mod tidy -diff
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: clean
clean:
	go clean

.PHONY: build
build:
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

.PHONY: run
run: build
	/tmp/bin/${binary_name}
