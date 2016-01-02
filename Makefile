GOROOT_BOOTSTRAP := $(abspath ../go_bootstrap)
GOROOT           := $(abspath ../go)

.PHONY: build
build: $(GOROOT)/bin/go update-stdlib

GODEPS = $(shell find $(GOROOT)/src/cmd $(GOROOT)/src/go -name "*.go")
$(GOROOT)/bin/go: $(GODEPS)
	cd $(GOROOT)/src && \
	  env GOROOT_BOOTSTRAP=$(GOROOT_BOOTSTRAP) CGO_ENABLED=0 ./make.bash

.PHONY: update-stdlib
update-stdlib:
	cd $(GOROOT) && git clean -q -df -- src/
	rsync -a src/ $(GOROOT)/src/

.PHONY: patch
patch:
	find patches/ -name "*.diff" -exec git apply --directory=$(GOROOT) {} \;

.PHONY: unpatch
unpatch:
	find patches/ -name "*.diff" -exec git apply -R --directory=$(GOROOT) {} \;

.PHONY:
clean:
	cd $(GOROOT) && git clean -dfx && git checkout .
