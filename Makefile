PROGNAME ?= philote
SOURCES = *.go src/**/*.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date
GPM ?= gpm

-include config.mk

$(PROGNAME):  bin $(SOURCES) $(DEPS) $(LUA_SOURCES) bin/philote-admin | $(dir $(PROGNAME))
	go build -o bin/$(PROGNAME)

bin/philote-admin: admin/*.go
	cd admin && go build -o ../$@

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
