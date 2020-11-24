# Go and compilation related variables
BUILD_DIR ?= out

BINARY_NAME := admin-helper
RELEASE_DIR ?= release

LDFLAGS := -extldflags='-static' -s -w

# Add default target
.PHONY: all
all: build

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -fr release

$(BUILD_DIR)/macos-amd64/$(BINARY_NAME):
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/macos-amd64/$(BINARY_NAME) ./main.go

$(BUILD_DIR)/linux-amd64/$(BINARY_NAME):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) ./main.go

$(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe ./main.go

.PHONY: cross ## Cross compiles all binaries
cross: $(BUILD_DIR)/macos-amd64/$(BINARY_NAME) $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe

.PHONY: release
release: clean cross
	mkdir $(RELEASE_DIR)
	tar cJSf $(RELEASE_DIR)/admin-helper-macos-amd64.tar.xz -C $(BUILD_DIR)/macos-amd64 $(BINARY_NAME)
	tar cJSf $(RELEASE_DIR)/admin-helper-linux-amd64.tar.xz -C $(BUILD_DIR)/linux-amd64 $(BINARY_NAME)
	tar cJSf $(RELEASE_DIR)/admin-helper-windows-amd64.tar.xz -C $(BUILD_DIR)/windows-amd64 $(BINARY_NAME).exe

	pushd $(RELEASE_DIR) && sha256sum * > sha256sum.txt && popd

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./main.go

