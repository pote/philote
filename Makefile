PROGNAME = philote
DEPS = main.go

$(PROGNAME): $(DEPS)
	go build

test: $(DEPS)
	go test

clean:
	rm $(PROGNAME)

.PHONY: test clean
