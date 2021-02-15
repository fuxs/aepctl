PLATFORM?=local
BIN=$(CURDIR)/bin
DIST=$(CURDIR)/dist
ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
	AEPCTL=$(BIN)/aepctl.exe
	AEPCTLDBG=$(BIN)/aepctldbg.exe
else
	AEPCTL=$(BIN)/aepctl
	AEPCTLDBG=$(BIN)/aepctldbg
endif

.PHONY: build build-dbg build-dist build-win-386 build-win-amd64 debug dependencies lint vet build-in-container lint-in-container

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
	go build -gcflags="all=-N -l" -o "$(AEPCTLDBG)"  main.go

debug: build-dbg
	dlv --listen=:2345 --headless --api-version=2 exec "$(AEPCTLDBG)" -- configure

build: dependencies $(bin) 
	go build -o "$(AEPCTL)" main.go

build-dist: build-win-386 build-win-amd64

build-win-386: dependencies
	@GOOS=windows GOARCH=386 go build -o $(DIST)/windows/386/bin/aepctl.exe main.go
	@tar -C $(DIST)/windows/386 -czf $(DIST)/aepctl-windows-386.tgz bin/aepctl.exe
	@shasum -a 256 $(DIST)/aepctl-windows-386.tgz | head -c 64 > $(DIST)/aepctl-windows-386.tgz.sha256

build-win-amd64: dependencies
	@GOOS=windows GOARCH=amd64 go build -o $(DIST)/windows/amd64/bin/aepctl.exe main.go
	@tar -C $(DIST)/windows/amd64 -czf $(DIST)/aepctl-windows-amd64.tgz bin/aepctl.exe
	@shasum -a 256 $(DIST)/aepctl-windows-amd64.tgz | head -c 64 > $(DIST)/aepctl-windows-amd64.tgz.sha256

dependencies:
	go get ./...

vet:
	go vet ./...

lint:
	golangci-lint run --timeout 10m0s ./...