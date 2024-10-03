APP_NAME := golang-simple-cache

BUILD_DIR := build

CMD_DIR := ./cmd

CONFIG_PATH ?= ./config/config.yaml

LDFLAGS := -ldflags "-s -w"

GO := go

all: build

build:
	@echo "==> Building the project..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: build
	@echo "==> Running the application..."
	./$(BUILD_DIR)/$(APP_NAME) -config $(CONFIG_PATH)

clean:
	@echo "==> Cleaning build directory..."
	rm -rf $(BUILD_DIR)
