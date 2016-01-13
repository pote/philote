PROGNAME ?= philote
SOURCES = *.go src/**/*.go
LUA_SOURCES = $(patsubst lua/%.lua,src/lua/scripts/%.go,$(wildcard lua/*.lua))
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date
GPM ?= gpm

$(PROGNAME):  $(SOURCES) $(DEPS) $(LUA_SOURCES) philote-cli | $(dir $(PROGNAME))
	go build -o $(PROGNAME)

philote-cli: cli/*.go
	cd cli && go build -o $@

server: $(PROGNAME)
	./$(PROGNAME)
test: $(PROGNAME) $(SOURCES)
	go test

clean:
	rm $(PROGNAME)

dependencies: $(DEPS)

$(DEPS): Godeps | $(dir $(DEPS))
	$(GPM) get
	touch $@

$(dir $(DEPS)):
	mkdir -p $@

$(dir $(PROGNAME)):
	mkdir -p $@

##
# Lua Scripts -> Go files
##

src/lua/scripts/%.go: src/lua/scripts lua/%.lua
	./script/asset_to_go "$*"

##
# Directories
##

src/lua/script:
	mkdir -p $@


##
# You're a PHONY! Just a big, fat PHONY.
##

.PHONY: run test clean dependencies
