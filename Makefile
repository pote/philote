PROGNAME ?= philote
DEPS = main.go

$(PROGNAME): $(DEPS) .dependencies/up-to-date
	mkdir -p $(@D)
	go build -o $@

run: $(PROGNAME)
	$(PROGNAME)

test: $(DEPS)
	go test

clean:
	rm $(PROGNAME)

.dependencies/up-to-date: Godeps
	gpm install
	touch $@

.PHONY: run test clean
