patterns = "/application|/outerlib"
PROJ_ROOT = $(shell pwd | awk -F$(patterns) '{print $$1}')

.PHONY: all clean

all:
	@echo
	go build -i

clean:
	@echo
	go clean -x
