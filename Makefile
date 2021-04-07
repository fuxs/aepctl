PLATFORM?=local
BIN=bin
GEN=gen
GEN_FILE=$(GEN)/pkged.go
DIST=dist

ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
	AEPCTL=$(BIN)/aepctl.exe
	AEPCTLDBG=$(BIN)/aepctldbg.exe
	DATE_CMD=powershell Get-Date -format "{dd-MMM-yyyy HH:mm}"
else
	AEPCTL=$(BIN)/aepctl
	AEPCTLDBG=$(BIN)/aepctldbg
	DATE_CMD=date
endif

COMMIT=$(shell git rev-list -1 HEAD)
BUILD_TIME=$(shell $(DATE_CMD))
VERSION=$(shell git describe --tags)
LDFLAGS=-X 'github.com/fuxs/aepctl/cmd/version.Commit=${COMMIT}'
LDFLAGS+=-X 'github.com/fuxs/aepctl/cmd/version.BuildTime=${BUILD_TIME}'
LDFLAGS+=-X 'github.com/fuxs/aepctl/cmd.Version=${VERSION}'

.PHONY: build build-dbg build-dist build-win-386 build-win-amd64 debug dependencies lint vet build-in-container lint-in-container

$(BIN):
	@mkdir $(BIN)

build-in-container: $(BIN)
	DOCKER_BUILDKIT=1 docker build . --target bin \
	--file aepctl.dockerfile \
	--platform ${PLATFORM} \
	--output bin/

lint-in-container:
	DOCKER_BUILDKIT=1 docker build . --target lint \
	--file aepctl.dockerfile

build-dbg: dependencies 
	go build -gcflags="all=-N -l" -o "$(AEPCTLDBG)"  main.go

debug: build-dbg
	dlv --listen=:2345 --headless --api-version=2 exec "$(AEPCTLDBG)" -- configure

build: dependencies
	@go build -ldflags="$(LDFLAGS)" -o "$(AEPCTL)" main.go

build-dist: build-win-386 build-win-amd64

build-win-386: dependencies
	@GOOS=windows GOARCH=386 go build -ldflags="$(LDFLAGS)" -o $(DIST)/windows/386/bin/aepctl.exe main.go
	@tar -C $(DIST)/windows/386 -czf $(DIST)/aepctl-windows-386.tgz bin/aepctl.exe
	@shasum -a 256 $(DIST)/aepctl-windows-386.tgz | head -c 64 > $(DIST)/aepctl-windows-386.tgz.sha256
	@shasum -a 256 $(DIST)/windows/386/bin/aepctl.exe | head -c 64 > $(DIST)/windows/386/bin/aepctl.exe.sha256

build-win-amd64: dependencies
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST)/windows/amd64/bin/aepctl.exe main.go
	@tar -C $(DIST)/windows/amd64 -czf $(DIST)/aepctl-windows-amd64.tgz bin/aepctl.exe
	@shasum -a 256 $(DIST)/aepctl-windows-amd64.tgz | head -c 64 > $(DIST)/aepctl-windows-amd64.tgz.sha256
	@shasum -a 256 $(DIST)/windows/amd64/bin/aepctl.exe | head -c 64 > $(DIST)/windows/amd64/bin/aepctl.exe.sha256

dependencies:
	@go get ./...

vet:
	go vet ./...

lint:
	golangci-lint run --timeout 10m0s ./...