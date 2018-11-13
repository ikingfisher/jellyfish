FILTER_OUT = $(wildcard makefile*) README.md make.include src pkg

ifndef SUBDIRS
	SUBDIRS = $(filter-out $(FILTER_OUT), $(wildcard *))
endif

.PHONY: all clean subdirs $(SUBDIRS)

all: $(SUBDIRS)

clean: $(SUBDIRS)
	-rm pkg -r

$(SUBDIRS):
	-$(MAKE) $(MAKECMDGOALS) -C $@