GOROOT_BOOTSTRAP := $(abspath build/go_bootstrap)
GOROOT           := $(abspath build/go)

.PHONY: build
build: $(GOROOT)/bin/go update-stdlib
	GOOS=atman $(GOROOT)/bin/go build -a runtime

GODEPS = $(shell find $(GOROOT)/src/cmd $(GOROOT)/src/go -name "*.go")
$(GOROOT)/bin/go: $(GODEPS)
	cd $(GOROOT)/src && \
	  env GOROOT_BOOTSTRAP=$(GOROOT_BOOTSTRAP) CGO_ENABLED=0 ./make.bash

.PHONY: update-stdlib
update-stdlib:
	cd $(GOROOT) && git clean -q -df -- src/
	rsync -a src/ $(GOROOT)/src/

PATCHES = $(shell find $(PWD)/patches -name "*.diff")
.PHONY: patch
patch:
	cd $(GOROOT) && git apply $(PATCHES)

.PHONY: unpatch
unpatch:
	cd $(GOROOT) && git apply -R $(PATCHES)

.PHONY: clean
clean:
	cd $(GOROOT) && git clean -dfx && git checkout .
