patterns = "/jellyfish|/outerlib"
PROJ_ROOT = $(shell pwd | awk -F$(patterns) '{print $$1}')
include $(PROJ_ROOT)/jellyfish/make.include

.PHONY: all clean

all:
	@echo
	go build -i

clean:
	@echo
	go clean -x
