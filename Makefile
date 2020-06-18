ALL := web_service spider
BUILD_TIME=$(shell date '+%Y-%m-%d %H:%M:%S')
BUILD_VERSION=$(shell git rev-parse HEAD)
BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GO_VERSION=$(shell go version)
BUILD_PATH=$(shell pwd)
WHO=$(shell git config user.name)
# LDFLAGS=-ldflags "-X 'common/version.BuildTime=$(BUILD_TIME)' -X 'common/version.BuildVersion=$(BUILD_VERSION)' -X 'common/version.BuildBranch=$(BUILD_BRANCH)' -X 'common/version.BuildName=$(WHO)' -X 'common/version.BuildPath=$(BUILD_PWD)' -X 'common/version.GoVersion=$(GO_VERSION)'"
FABBUILD := GOOS=linux GOARCH=amd64 CGO_ENABLE=0 go build
LOCALBUILD := go build
SSHKEY="~/Downloads/tmp"
all: $(ALL)
	@echo "build $@ over!"

fab_all: fab_spider fab_web_service

web_service:
	$(LOCALBUILD) -o bin/web_service src/web_service/main.go
	@echo "debug build $@ over!"
	$(FABBUILD) -o bin/linux/web_service src/web_service/main.go
	@echo "release build $@ over!"
spider:
	$(LOCALBUILD) -o bin/spider src/spider/main.go
	@echo "debug build $@ over!"
	$(FABBUILD) -o bin/linux/spider src/spider/main.go
	@echo "release build $@ over!"

fab_spider:
	scp -i $(SSHKEY) bin/linux/spider root@myapig.info:/root
fab_web_service:
	scp -i $(SSHKEY) bin/linux/web_service root@myapig.info:/root
clean:
	rm -rf bin/*
