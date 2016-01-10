PROGNAME ?= philote
SOURCES = *.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date
GPM ?= gpm

$(PROGNAME): $(SOURCES) $(DEPS) | $(dir $(PROGNAME))
	go build -o $(PROGNAME)

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

.PHONY: run test clean dependencies deploy provision ansible/files/philote
