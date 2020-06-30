.PHONY: lint test all check-counterfeiter gen-mocks release

all: test

OUTPUT = ./bin/prel
GO_SOURCES = $(shell find . -type f -name '*.go')
VERSION ?= $(shell cat VERSION)
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")
LDFLAGS_VERSION = -X github.com/pivotal/scdf-k8s-prel/pkg/commands.cliVersion=$(VERSION) \
				  -X github.com/pivotal/scdf-k8s-prel/pkg/commands.cliGitsha=$(GITSHA) \
				  -X github.com/pivotal/scdf-k8s-prel/pkg/commands.cliGitdirty=$(GITDIRTY)

test:
	GO111MODULE=on go test ./pkg/...

lint:
	./scripts/check-lint.sh

check-counterfeiter:
    # Use go get in GOPATH mode to install/update counterfeiter. This avoids polluting go.mod/go.sum.
	@which counterfeiter > /dev/null || (echo counterfeiter not found: issue "GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter" && false)

gen-mocks: check-counterfeiter

prel: $(GO_SOURCES)
	GO111MODULE=on go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT) cmd/prel/main.go

release: $(GO_SOURCES) test
	GOOS=darwin   GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     cmd/prel/main.go && tar -czf irel-darwin-amd64.tgz  $(OUTPUT)     && rm -f $(OUTPUT)
	GOOS=linux    GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     cmd/prel/main.go && tar -czf irel-linux-amd64.tgz   $(OUTPUT)     && rm -f $(OUTPUT)
	GOOS=windows  GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT).exe cmd/prel/main.go && zip -mq  irel-windows-amd64.zip $(OUTPUT).exe && rm -f $(OUTPUT).exe
