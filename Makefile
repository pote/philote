PROGNAME ?= philote
SOURCES = *.go src/**/*.go
LUA_SOURCES = $(patsubst lua/%.lua,src/lua/scripts/%.go,$(wildcard lua/*.lua))
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
	go test

clean:
	rm -rf pkg/

dependencies: $(DEPS)

cross-compile: clean
	script/cross-compile

config.mk:
	@./configure

install: philote bin/philote-admin
	install -d $(prefix)/bin
	install -m 0755 bin/philote /usr/local/bin
	install -m 0755 bin/philote-admin /usr/local/bin

uninstall:
	rm -f $(prefix)/bin/philote
	rm -f $(prefix)/bin/philote-admin

$(DEPS): Godeps | $(dir $(DEPS))
	$(GPM) get
	touch $@

##
# Lua Scripts -> Go files
##

src/lua/scripts/%.go: src/lua/scripts lua/%.lua
	./script/asset_to_go "$*"

##
# Directories
##

src/lua/scripts bin:
	mkdir -p $@

$(dir $(DEPS)):
	mkdir -p $@

$(dir $(PROGNAME)):
	mkdir -p $@

##
# You're a PHONY! Just a big, fat PHONY.
##

.PHONY: run test clean dependencies cross-compile
