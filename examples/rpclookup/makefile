FILTER_OUT = $(wildcard Makefile*) inc lib makefile

ifndef SUBDIRS
	SUBDIRS = $(filter-out $(FILTER_OUT), $(wildcard *))
endif

.PHONY: all clean subdirs $(SUBDIRS)

all: $(SUBDIRS)

clean: $(SUBDIRS)

$(SUBDIRS):
	-$(MAKE) $(MAKECMDGOALS) -C $@