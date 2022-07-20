-include eng/Makefile

.DEFAULT_GOAL = build
.PHONY: \
	lint \
	install \
	-install-%

BUILD_VERSION=$(shell git rev-parse --short HEAD)
GO_LDFLAGS=-X 'github.com/Carbonfrost/autogun/pkg/internal/build.Version=$(BUILD_VERSION)'

lint:
	$(Q) go run honnef.co/go/tools/cmd/staticcheck -checks 'all,-ST*' $(shell go list ./...)

install: -install-autogun

-install-%: build -check-env-PREFIX -check-env-_GO_OUTPUT_DIR
	$(Q) eng/install "${_GO_OUTPUT_DIR}/$*" $(PREFIX)/bin
