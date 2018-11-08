patterns = "/application|/outerlib"
PROJ_ROOT = $(shell pwd | awk -F$(patterns) '{print $$1}')

GOPATH := $(PROJ_ROOT)/outerlib/go/FrameLibs
GOPATH := $(GOPATH):$(PROJ_ROOT)/application/shuoshuo/go
GOPATH := $(GOPATH):$(PROJ_ROOT)/outerlib/go/ServiceCode
GOPATH := $(GOPATH):$(PROJ_ROOT)/outerlib/go/ExtendSTL
GOPATH := $(GOPATH):$(PROJ_ROOT)/outerlib/go/GitHub
GOPATH := $(GOPATH):$(shell pwd)

.PHONY: all clean

all:
	@echo
	go build -i

clean:
	@echo
	go clean -x