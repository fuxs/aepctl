PLATFORM?=local
BIN=$(CURDIR)/bin
AEPCTL=$(BIN)/aepctl

.PHONY: build dependencies vet build-in-container lint-in-container

$(BIN):
	@echo -e "\033[1;32mCreating new bin directory\033[0m"
	mkdir -m 755 $(BIN)

build-in-container: $(BIN)
	DOCKER_BUILDKIT=1 docker build . --target bin \
	--file aepctl.dockerfile \
	--platform ${PLATFORM} \
	--output bin/

lint-in-container:
	DOCKER_BUILDKIT=1 docker build . --target lint \
	--file aepctl.dockerfile

build: dependencies $(bin) 
	go build -o bin/aepctl main.go

dependencies:
	go get ./...

vet:
	go vet ./...