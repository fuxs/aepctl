# syntax = docker/dockerfile:1-experimental

FROM --platform=${BUILDPLATFORM} golang:1.15.5-alpine AS base
WORKDIR /go/src/app
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
#COPY . . # replaced by --mount=target=.

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/aepctl main.go

FROM base AS unit-test
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  go test -v .

FROM golangci/golangci-lint:v1.31-alpine AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/.cache/golangci-lint \
  golangci-lint run --timeout 10m0s ./...

FROM scratch AS bin-unix
COPY --from=build /out/aepctl /
ENTRYPOINT [ "/aepctl" ]

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /out/aepctl /aepctl.exe
ENTRYPOINT [ "/aepctl.exe" ]

FROM bin-${TARGETOS} AS bin