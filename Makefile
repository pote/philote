PROGNAME ?= philote
SOURCES = *.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date
GPM ?= gpm

-include config.mk

all:
	go build -o bin/philote && bin/philote

$(PROGNAME):  bin $(SOURCES) $(DEPS) | $(dir $(PROGNAME))
	go build -o bin/$(PROGNAME)

server: $(PROGNAME)
	./bin/$(PROGNAME)

test: $(PROGNAME) $(SOURCES)
	LOG=error go test

clean:
	rm -rf pkg/

dependencies: $(DEPS)

cross-compile: clean
	script/cross-compile

config.mk:
	@./configure

install: philote
	install -d $(prefix)/bin
	install -m 0755 bin/philote /usr/local/bin

uninstall:
	rm -f $(prefix)/bin/philote

$(DEPS): Godeps | $(dir $(DEPS))
	$(GPM) get
	touch $@

##
# Directories
##

$(dir $(PROGNAME)) $(dir $(DEPS)) bin:
	mkdir -p $@


##
# You're a PHONY! Just a big, fat PHONY.
##

.PHONY: run test clean dependencies cross-compile
