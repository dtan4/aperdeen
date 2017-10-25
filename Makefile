NAME      := aperdeen
VERSION   := v0.1.0
REVISION  := $(shell git rev-parse --short HEAD)

SRCS      := $(shell find . -name '*.go' -type f)
LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""
PKG_DIRS  := $(shell go list ./... | grep -v vendor | grep -v apimock)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	echo "" > coverage.txt
	set -e; \
	for d in $(PKG_DIRS); do \
		go test -coverprofile=profile.out -covermode=atomic -v $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build:
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep mockgen
	dep ensure -v

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: generate
generate: mockgen
	go generate $(PKG_DIRS)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: mockgen
mockgen:
ifeq ($(shell command -v mockgen 2> /dev/null),)
	go get -v github.com/golang/mock/mockgen
endif

.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)

.PHONY: test
test:
	go test -cover -v $(PKG_DIRS)

.PHONY: update-deps
update-deps: dep
	dep ensure -update -v
