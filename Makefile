PLATFORM?=local
BIN=$(CURDIR)/bin
AEPCTL=$(BIN)/aepctl

.PHONY: build debug dependencies lint vet build-in-container lint-in-container

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

build-dbg: dependencies $(bin) 
	go build -gcflags="all=-N -l" -o bin/aepctldbg main.go

debug: build-dbg
	dlv --listen=:2345 --headless --api-version=2 exec bin/aepctldbg -- configure

build: dependencies $(bin) 
	go build -o bin/aepctl main.go

dependencies:
	go get ./...

vet:
	go vet ./...

lint:
	golangci-lint run --timeout 10m0s ./...