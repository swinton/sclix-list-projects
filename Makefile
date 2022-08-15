# Build system related variables
GOCMD = go
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

# Project related variables
WORKING_DIR = $(CURDIR)
BUILD_DIR = $(WORKING_DIR)/_bin

# build related variables
_SEPARATOR = _
_EXE_POSTFIX = 

ifeq ($(GOOS),windows)
	_EXE_POSTFIX = .exe
endif

# some globally assembled variables
APPLICATION_NAME = $(shell basename -s .git `git config --get remote.origin.url`)
PLATFORM_STRING = $(GOOS)$(_SEPARATOR)$(GOARCH)
EXECUTABLE_NAME = $(APPLICATION_NAME)$(_SEPARATOR)$(PLATFORM_STRING)$(_EXE_POSTFIX)

# Make directories per convention
prefix = /usr/local
exec_prefix = $(prefix)
bindir = $(exec_prefix)

# some make file variables
LOG_PREFIX = --

.PHONY: clean
clean:
	@$(GOCMD) clean
	@rm -f -r $(BUILD_DIR)

.PHONY: format
format:
	gofmt -w -l -e .

.PHONY: lint
lint:
	./scripts/lint.sh

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

SRCS = $(shell find $(WORKING_DIR) -type f -name '*.go')

$(BUILD_DIR)/$(EXECUTABLE_NAME): $(BUILD_DIR) $(SRCS)
	@echo "Building..."
	@echo "$(LOG_PREFIX) Building ( $(BUILD_DIR)/$(EXECUTABLE_NAME) )"
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build -o $(BUILD_DIR)/$(EXECUTABLE_NAME) $(WORKING_DIR)/cmd/main.go

.PHONY: build
build: $(BUILD_DIR)/$(EXECUTABLE_NAME)
# suppress "Nothing to be done for `build'."
	@echo > /dev/null 

.PHONY: test
test: 
	@echo "$(LOG_PREFIX) Testing"
	@$(GOCMD) test -cover ./...

$(DESTDIR)$(bindir):
	@mkdir -p $(DESTDIR)$(bindir)

.PHONY: install
install: $(DESTDIR)$(bindir) build
	@echo "$(LOG_PREFIX) Installing $(EXECUTABLE_NAME) ( $(DESTDIR)$(bindir) )"
	@cp -f $(BUILD_DIR)/$(EXECUTABLE_NAME) $(DESTDIR)$(bindir)
	@cp -f ./extension.json $(DESTDIR)$(bindir)

.PHONY: help
help:
	@echo "Main targets:"
	@echo "$(LOG_PREFIX) format"
	@echo "$(LOG_PREFIX) lint"
	@echo "$(LOG_PREFIX) build"
	@echo "$(LOG_PREFIX) test"
	@echo "$(LOG_PREFIX) clean"
	@echo "\nAvailable parameter:"
	@echo "$(LOG_PREFIX) GOOS                       Specify Operating System to compile for (see golang GOOS, default=$(GOOS))"
	@echo "$(LOG_PREFIX) GOARCH                     Specify Architecture to compile for (see golang GOARCH, default=$(GOARCH))"
