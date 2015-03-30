PROGNAME ?= philote
SOURCES = main.go access_token.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date

$(PROGNAME): $(SOURCES) $(DEPS)
	mkdir -p $(@D)
	go build -o $@

run: $(PROGNAME)
	$(PROGNAME)

test: $(PROGNAME) $(SOURCES)
	go test

clean:
	rm $(PROGNAME)

dependencies: $(DEPS)

$(DEPS): Godeps | $(dir $(DEPS))
	gpm install
	touch $@

$(dir $(DEPS)):
	mkdir -p $@

.PHONY: run test clean dependencies
