GO := GOPATH=$(shell go env GOROOT)/bin:"$(shell pwd)" GOOS=$(OS) GOARCH=$(ARCH) go
#GO := GOPATH=$(shell go env GOROOT)/bin:"$(shell pwd)" go
GOGETTER := GOPATH="$(shell pwd)" GOOS=$(OS) GOARCH=$(ARCH) go get -u
.PHONY: all scheduler

all: clean go_get_deps scheduler

scheduler:
	$(GO) build $(GOOPTS) -o bin/scheduler scheduler.go scheduler_funcs.go util.go db.go restclient.go log.go state.go commands.go config.go cron_eval.go

go_get_deps:
	$(GOGETTER) github.com/Sirupsen/logrus
	$(GOGETTER) github.com/gorhill/cronexpr
	$(GOGETTER) github.com/jmcvetta/napping
	$(GOGETTER) github.com/mattn/go-sqlite3

clean:
	rm -rf bin src/github.com

