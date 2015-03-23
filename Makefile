PROGNAME ?= philote
DEPS = main.go

$(PROGNAME): $(DEPS)
	mkdir -p $(@D)
	go build -o $@

run: $(PROGNAME)
	$(PROGNAME)

test: $(DEPS)
	go test

clean:
	rm $(PROGNAME)

.PHONY: run test clean
