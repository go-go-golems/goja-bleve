.PHONY: gifs logcopter-generate logcopter-check test-vectors xgoja-build-vectors xgoja-smoke-vectors

all: gifs

VERSION=v0.1.14
GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target
XGOJA_VERSION ?= v0.8.3
XGOJA_VECTOR_SPEC ?= xgoja-vectors.yaml
XGOJA_VECTOR_WORK_DIR ?= /tmp/goja-bleve-vector-work

TAPES=$(wildcard doc/vhs/*tape)
gifs: $(TAPES)
	for i in $(TAPES); do vhs < $$i; done

docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

lint:
	GOWORK=off golangci-lint run -v

lintmax:
	GOWORK=off golangci-lint run -v --max-same-issues=100

gosec:
	GOWORK=off go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude-generated -exclude=G101,G304,G301,G306 -exclude-dir=.history ./...

govulncheck:
	GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	GOWORK=off go test ./...

test-vectors:
	GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1

xgoja-build-vectors:
	cd cmd/goja-bleve && GOWORK=off go run github.com/go-go-golems/go-go-goja/cmd/xgoja@$(XGOJA_VERSION) build \
		-f $(XGOJA_VECTOR_SPEC) \
		--work-dir $(XGOJA_VECTOR_WORK_DIR) \
		--keep-work \
		--xgoja-version $(XGOJA_VERSION)

xgoja-smoke-vectors: xgoja-build-vectors
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector knn --output json
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector hybrid --output json

build:
	GOWORK=off go generate ./...
	GOWORK=off go build ./...

logcopter-generate:
	GOWORK=off go generate ./...

logcopter-check:
	GOWORK=off go tool logcopter-gen -area-prefix go-go-golems.goja_bleve -strip-prefix github.com/go-go-golems/goja-bleve -check ./pkg/...

goreleaser:
	GOWORK=off goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $(shell svu major)

tag-minor:
	git tag $(shell svu minor)

tag-patch:
	git tag $(shell svu patch)

release:
	git push origin --tags
	GOWORK=off GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/goja-bleve@$(shell svu current)

bump-go-go-golems:
	@deps="$$(awk '/^require[[:space:]]+github\.com\/go-go-golems\// { print $$2 } /^[[:space:]]*github\.com\/go-go-golems\// { print $$1 }' go.mod | sort -u)"; \
	if [ -z "$$deps" ]; then \
		echo "No github.com/go-go-golems dependencies in go.mod"; \
	else \
		echo "Bumping go-go-golems dependencies:"; \
		echo "$$deps"; \
		for dep in $$deps; do GOWORK=off go get "$${dep}@latest"; done; \
	fi
	GOWORK=off go mod tidy

GOJA_BLEVE_BINARY=$(shell which goja-bleve)
install:
	cd cmd/goja-bleve && GOWORK=off go build -o ../../dist/goja-bleve . && \
		cp ../../dist/goja-bleve $(GOJA_BLEVE_BINARY)
