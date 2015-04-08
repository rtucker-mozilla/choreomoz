GO := GOPATH=$(shell go env GOROOT)/bin:"$(shell pwd)" GOOS=$(OS) GOARCH=$(ARCH) go
#GO := GOPATH=$(shell go env GOROOT)/bin:"$(shell pwd)" go
GOGETTER := GOPATH="$(shell pwd)" GOOS=$(OS) GOARCH=$(ARCH) go get -u
.PHONY: all chorizo

all: clean go_get_deps chorizo

test_cron_eval:
	$(GO) test cron_eval_test.go cron_eval.go

test_util:
	$(GO) test src/libchorizo/util/util_test.go src/libchorizo/util/util.go

chorizo:
	$(GO) build $(GOOPTS) -o bin/chorizo chorizo.go chorizo_funcs.go commands.go cron_eval.go

go_get_deps:
	$(GOGETTER) github.com/Sirupsen/logrus
	$(GOGETTER) github.com/gorhill/cronexpr
	$(GOGETTER) github.com/jmcvetta/napping
	$(GOGETTER) github.com/mattn/go-sqlite3
	$(GOGETTER) code.google.com/p/gcfg

clean:
	rm -rf bin src/github.com src/bitbucket.org src/code.google.com

