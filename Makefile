PLATFORM?=local
BIN=$(CURDIR)/bin
AEPCTL=$(BIN)/aepctl

$(info go filesz $(GO_FILES))

.PHONY: build

$(BIN):
	@echo -e "\033[1;32mCreating new bin directory\033[0m"
	mkdir -m 755 $(BIN)

build: $(BIN)
	DOCKER_BUILDKIT=1 docker build . --target bin \
	--file aepctl.dockerfile \
	--platform ${PLATFORM} \
	--output bin/

lint:
	DOCKER_BUILDKIT=1 docker build . --target lint \
	--file aepctl.dockerfile