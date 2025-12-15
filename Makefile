.PHONY: all build clean test run-api

all: build

build:
	@echo "building m8s components..."
	@mkdir -p bin
	@go build -o bin/apiserver ./cmd/apiserver
	@go build -o bin/scheduler ./cmd/scheduler
	@go build -o bin/kubelet ./cmd/kubelet
	@go build -o bin/m8s ./cmd/m8s
	@echo "build complete! binarier in ./bin/"

test:
	@echo "running tests..."
	@go test -v ./...

run-api:
	@./bin/apiserver

clean:
	@echo "cleaning"
	@rm -rf bin/
	@rm -f /var/lib/m8s/state.json
	@echo "cleaned!"

install: build
	@echo "installing m8s binaries..."
	@sudo cp bin/* /usr/local/bin/
	@echo "installed!"
