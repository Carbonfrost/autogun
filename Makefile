# Copyright 2025, 2026 The Autogun Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

-include eng/Makefile

.DEFAULT_GOAL = build
.PHONY: \
	lint \
	install \
	-install-%

GO_LDFLAGS=

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

json_info:
	@ go run -tags json_marshal ./cmd/autogun > docs/autogun.json_info.json