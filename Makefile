APP_NAME=docshell
BUILD_DIR=bin
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build run lint clean

all: build run

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

run: build
	@echo "Running $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

lint:
	@echo "Running gofmt and go vet..."
	@gofmt -l -s $(GO_FILES)
	@go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

docker-build:
	docker build -t docshell .

docker-run:
	docker run --rm -d docshell

docker-rund:
	docker run --rm -it -p 8080:8080 docshell