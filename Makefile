SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sdees

VERSION=1.0.1
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go get github.com/maxwellhealth/go-gpg
	go get github.com/pkg/sftp
	go get github.com/mitchellh/go-homedir
	go get github.com/urfave/cli
	go get github.com/jcelliott/lumber
	go build ${LDFLAGS} -o ${BINARY} main.go

# .PHONY: install
# install:
# 	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
