# Get relative paths of all "main" packages
MAIN_PACKAGE ?= $(shell go list -f '{{if (eq .Name "main")}}.{{slice .ImportPath (len .Module.Path)}}{{end}}' ./...)

TEST_PACKAGE ?= ./...

.PHONY: build
build:
	go build -gcflags '-e' ./...

.PHONY: test
test:
	go test $(TEST_PACKAGE)

.PHONY: race
race:
	go test -race $(TEST_PACKAGE)

.PHONY: bench
bench:
	go test -bench $(TEST_PACKAGE)

.PHONY: tags
tags:
	gotags -f tags -R .

.PHONY: cover
cover:
	mkdir -p tmp
	go test -coverprofile tmp/_cover.out $(TEST_PACKAGE)
	go tool cover -html tmp/_cover.out -o tmp/cover.html

.PHONY: checkall
checkall: vet staticcheck

.PHONY: vet
vet:
	go vet $(TEST_PACKAGE)

.PHONY: staticcheck
staticcheck:
	staticcheck $(TEST_PACKAGE)

.PHONY: clean
clean:
	go clean
	rm -f tags
	rm -f tmp/_cover.out tmp/cover.html

list-upgradable-modules:
	@go list -m -u -f '{{if .Update}}{{.Path}} {{.Version}} [{{.Update.Version}}]{{end}}' all

# Build all "main" packages
.PHONY: main-build
main-build:
	@for d in $(MAIN_PACKAGE) ; do \
	  echo "cd $$d && go build -gcflags '-e'" ; \
	  ( cd $$d && go build -gcflags '-e' ) ; \
	done

# Clean all "main" packages
.PHONY: main-clean
main-clean:
	@for d in $(MAIN_PACKAGE) ; do \
	  echo "cd $$d && go clean" ; \
	  ( cd $$d && go clean ) ; \
	done

# based on: github.com/koron-go/_skeleton/Makefile
# $Hash:46e18120a853c06f4546e24b5c4ba7e823abefe1f601db9b2dc289f5$
