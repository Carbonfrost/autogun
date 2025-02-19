-include eng/Makefile

.DEFAULT_GOAL = build
.PHONY: \
	lint \
	install \
	-install-%

BUILD_VERSION=$(shell git rev-parse --short HEAD)
GO_LDFLAGS=-X 'github.com/Carbonfrost/autogun/pkg/internal/build.Version=$(BUILD_VERSION)'

lint:
	$(Q) go tool gocritic check ./... 2>&1 || true
	$(Q) go tool revive ./... 2>&1 || true
	$(Q) go tool staticcheck -checks 'all,-ST*' $(shell go list ./...) 2>&1	

fmt: -fmt-hcl

-fmt-hcl:	
	$(Q) go tool hclfmt -w pkg/config/testdata/valid-examples/*.autog

install: -install-autogun

-install-%: build -check-env-PREFIX -check-env-_GO_OUTPUT_DIR
	$(Q) eng/install "${_GO_OUTPUT_DIR}/$*" $(PREFIX)/bin
