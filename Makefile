PROGNAME ?= philote
SOURCES = main.go access_token.go socket.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date

$(PROGNAME): $(SOURCES) $(DEPS) | $(dir $(PROGNAME))
	go build -o $(PROGNAME)

run: $(PROGNAME)
	./$(PROGNAME)

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

$(dir $(PROGNAME)):
	mkdir -p $@

.PHONY: run test clean dependencies
