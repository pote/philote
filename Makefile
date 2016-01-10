PROGNAME ?= philote
SOURCES = *.go
DEPS = $(firstword $(subst :, ,$(GOPATH)))/up-to-date
GPM ?= gpm

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
	$(GPM) get
	touch $@

$(dir $(DEPS)):
	mkdir -p $@

$(dir $(PROGNAME)):
	mkdir -p $@

##
# Provisioning and Deploy
##

provision:
	ansible-playbook -i ansible/inventory ansible/provision.yml

deploy: ansible/files/philote
	ansible-playbook -i ansible/inventory ansible/deploy.yml

ansible/files/philote:
	GOOS=linux GOARCH=amd64 go build -o $@

.PHONY: run test clean dependencies deploy provision ansible/files/philote
