.DEFAULT_GOAL := container

define ASCISKOPOS

      _               _ _ 
     | |             | (_)
  ___| | ________ ___| |_ 
 |_  / |/ /______/ __| | |
  / /|   <      | (__| | |
 /___|_|\_\      \___|_|_|
                          
                          

endef

export ASCISKOPOS

# http://misc.flogisoft.com/bash/tip_colors_and_formatting

RED=\033[0;31m
GREEN=\033[0;32m
ORNG=\033[38;5;214m
BLUE=\033[38;5;81m
NC=\033[0m

export RED
export GREEN
export NC
export ORNG
export BLUE


PUBLISH_TAG?=docker-ethos-core-univ-release.dr-uw2.adobeitc.com/ethos

help:
	@printf "\033[1m$$ASCISKOPOS $$NC\n"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | sort | xargs -n 1 -IXXX printf "\033[1mXXX $$NC\n"


build-container: ## builds and tags to the current VERSION
build-container: container tag-container
	@echo "Building zk-cli container ..."

tag-container: ## Apply docker upstream tags
tag-container:
	docker tag adobeplatform/zk-cli:`git rev-parse HEAD` $(PUBLISH_TAG)/zk-cli:`cat VERSION`; \
	docker tag adobeplatform/zk-cli:`git rev-parse HEAD` adobeplatform/zk-cli:latest 

upload-container: ## uploads to adobeplatform.  You need to have credentials.  Make sure you set DOCKER_CONFIG=`cd $$HOME/.docker-hub-f4tq/;pwd`
upload-container: build-container
	docker push $(PUBLISH_TAG)/zk-cli:`cat VERSION` 


strip:
	strip zk-cli

container:    ## makes zk-cli container 
container: Dockerfile  compile  strip
	@if [ ! -e /.dockerenv -o ! -z "$JENKINS_URL" ]; then \
		echo ; \
		echo ; \
		echo "------------------------------------------------" ; \
		echo "$@: Building zk-cli container image..." ; \
		echo "------------------------------------------------" ; \
		echo ; \
		docker build -f Dockerfile -t adobeplatform/zk-cli:`git rev-parse HEAD` . ; \
	else \
		@echo ; \
		echo "------------------------------------------------" ; \
		echo "$@: Running in Docker so skipping..." ; \
		echo "------------------------------------------------" ; \
		echo ; \
		env ; \
		echo ; \
	fi

######### Go ##########
.PHONY: ci test compile

SOURCES:=$(shell find . \( -name vendor  \) -prune -o  -name '*.go')
RESOURCES:=$(shell find ./resources -type f  | grep -v resources/bindata.go )


glide.lock: glide.yaml
# 	docker sometimes has trouble re-arranging ./vendor.  Better to just blow it away.  glide caches previous runs
	@glide --no-color install

install-deps:  ##  install dependencies.  Not usually needed outside of a container
install-deps: glide.lock 



install-tools:  ##  installs glide, golint, ginkgo, gomock, gomegs
install-tools:
	@which golint || go get -u github.com/golang/lint/golint
	@which cover || go get golang.org/x/tools/cmd/cover
	@test -d $$GOPATH/github.com/go-ini/ini || go get github.com/go-ini/ini
	@test -d $$GOPATH/github.com/jmespath/go-jmespath ||  go get github.com/jmespath/go-jmespath
	@which ginkgo || go get github.com/onsi/ginkgo/ginkgo
	@which gomega || go get github.com/onsi/gomega
	@which gomock || go get github.com/golang/mock/gomock
	@which mockgen || go get github.com/golang/mock/mockgen
	@which glide || go get github.com/Masterminds/glide
	@which go-bindata || go get -u github.com/jteeuwen/go-bindata/...

docker_resources: resources/bindata.go

resources/bindata.go: $(RESOURCES)
	go-bindata -pkg resources -o resources/bindata.go resources/...

zk-cli-linux-amd64:  ##   builds zk-cli server
zk-cli-linux-amd64: $(SOURCES)
	@if [ "$$(uname)" == "Linux" ]; then  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitSha=`git rev-parse HEAD` -X main.Version=`cat VERSION` -X main.BuildDate=`date +%s`" -o zk-cli-linux-amd64 ./ ;  else echo "Skipping Linux CGO build" ; fi


docker_compile:  install-deps docker_resources zk-cli-linux-amd64


docker_build: install-deps docker_compile


build: test container
	if [ "x$$sha" = "x" ] ; then sha=`git rev-parse HEAD`; fi ;\
	docker push adobeplatform/zk-cli:$$sha ;\
	docker push adobeplatform/zk-cli:latest

docker_lint: install-deps
	go tool vet -all app api config shared controllers models view server
	@DIRS=" shared/... " && FAILED="false" && \
	echo "gofmt -l *.go app api config controllers models shared  view server" && \
	GOFMT=$$(gofmt -l *.go  api app config shared controllers models view server ) && \
	if [ ! -z "$$GOFMT" ]; then echo -e "\nThe following files did not pass a 'go fmt' check:\n$$GOFMT\n" && FAILED="true"; fi; \
	for codeDir in $$DIRS; do \
		LINT="$$(golint $$codeDir)" && \
		if [ ! -z "$$LINT" ]; then echo "$$LINT" && FAILED="true"; fi; \
	done && \
	if [ "$$FAILED" = "true" ]; then exit 1; else echo "ok" ;fi


docker_test: install-deps docker_lint docker_compile
	go test -v --cover  $$(go list ./... | grep -v /vendor/)

docker_test_ci: install-deps docker_lint docker_compile
	go test -v --cover --timeout 60s $$(go list ./... | grep -v /vendor/)


docker_ci: docker_test_ci docker_compile
resources: ##  builds resources inside a docker container but leaves them in your directory
compile:  ##  compiles your project.  uses the dev-container
lint:  ##  lints the project.  Inside a container
ci:  ##  target for jenkins.  Inside a container 
test:  ##  tests the project.  Inside a container

resources compile lint test ci : dev-container
#   either ssh key or agent is needed to pull adobe-platform sources from git
#   this supplies to methods
#
	@SSH1="" ; SSH2="" ;\
	if [ "x$$sha" = "x" ] ; then sha=`git rev-parse HEAD`; fi ;\
        if [ ! -z "$$SSH_AUTH_SOCK" ] ; then SSH1="-e SSH_AUTH_SOCK=/root/.foo -v $$SSH_AUTH_SOCK:/root/.foo" ; fi ; \
        if [ -e $$HOME/.ssh/id_rsa ]; then SSH2="-v $$HOME/.ssh/id_rsa:/root/.ssh/id_rsa" ; fi ; \
	if [ ! -e /.dockerenv -o ! -z "$JENKINS_URL" ];  then \
	AWS=$$(env | grep AWS | xargs -n 1 -IXX echo -n ' -e XX') ;\
	echo ; \
	echo ; \
	echo "------------------------------------------------" ; \
	echo "Running target \"$@\" inside Docker container..." ; \
	echo "------------------------------------------------" ; \
	echo ; \
	docker run -i --rm $$SSH1 $$SSH2 $$AWS\
		--name=skopos_web_gate_make_docker_$@ \
		-e sha=$$sha \
        -v $$(pwd):/go/src/github.com/adobe-platform/zk-cli \
        -w /go/src/github.com/adobe-platform/zk-cli \
		adobe-platform/zk-cli:dev \
		make docker_$@ ;\
	else \
		make docker_$@ ;\
	fi


dev-container:  ##  makes dev-container.  runs make install-tools in dev-container.  Builds adobe-platform/zk-cli:dev
dev-container: Dockerfile-dev 
	@printf "\033[1m$$ASCISKOPOS $$NC\n"

	@if [ ! -e /.dockerenv -o ! -z "$JENKINS_URL" ]; then \
		echo ; \
		echo ; \
		echo "------------------------------------------------" ; \
		echo "$@: Building dev container image..." ; \
		echo "------------------------------------------------" ; \
		echo ; \
		docker images | grep 'adobe-platform/zk-cli' | awk '{print $$2}' | grep -q -E '^dev$$' ; \
		if [ $$? -ne 0 ]; then  \
			docker build -f Dockerfile-dev -t adobe-platform/zk-cli:dev . ; \
		fi ; \
	else \
		echo ; \
		echo "------------------------------------------------" ; \
		echo "$@: Running in Docker so skipping..." ; \
		echo "------------------------------------------------" ; \
		echo ; \
		env ; \
		echo ; \
	fi

clean: clean-dev

clean-dev:  ##  Remove the adobe-platform/zk-cli:dev
clean-dev:
	@if [ ! -e /.dockerenv -o ! -z "$JENKINS_URL" ]; then \
		if $$(docker ps | grep -q "adobe-platform/zk-cli:dev"); then \
			echo "You have a running dev container.  Stop it first before using clean-dev" ;\
			exit 10; \
		fi ; \
		docker images | grep 'adobe-platform/zk-cli' | awk '{print $$2}' | grep -q -E '^dev$$' ; \
		if [ $$? -eq 0 ]; then  \
			docker rmi adobe-platform/zk-cli:dev  ; \
		else \
			echo "No dev image" ;\
		fi ; \
	else \
		echo ; \
		echo "------------------------------------------------" ; \
		echo "$@: Running in Docker so skipping..." ; \
		echo "------------------------------------------------" ; \
		echo ; \
		env ; \
		echo ; \
	fi

run-dev:  ##  Runs the adobe-platform/zk-cli:dev container mounting the current directly.  Gives full dev environment.  Maps in your ssh-agent and keeps a bash-history outside the container so you have history between invocations.
run-dev: dev-container
#       save bash history in-between runs...
	@if [ ! -f $$HOME/.bash_history-zk-cli-dev ]; then touch $$HOME/.bash_history-zk-cli-dev; fi
#       mount the current directory into the dev build
#       map ssh-agent's auth-sock into the container instance.  the pipe needs to be on non-external volume hence /root/.foo
	@SSH1="" ; SSH2="" ;\
        if [ ! -z "$$SSH_AUTH_SOCK" ] ; then SSH1="-e SSH_AUTH_SOCK=/root/.foo -v $$SSH_AUTH_SOCK:/root/.foo" ; fi ; \
        if [ -e $$HOME/.ssh/id_rsa ]; then SSH2="-v $$HOME/.ssh/id_rsa:/root/.ssh/id_rsa" ; fi ; \
        AWS=$$(env | grep AWS | xargs -n 1 -IXX echo -n ' -e XX'); \
	docker run -i --rm --net host  $$SSH1 $$SSH2 $$AWS -e HISTSIZE=100000  -v $$HOME/.bash_history-zk-cli-dev:/root/.bash_history -v `pwd`:/go/src/github.com/adobe-platform/zk-cli -w /go/src/github.com/adobe-platform/zk-cli -t adobe-platform/zk-cli:dev bash ; \
	if [ $$? -ne 0 ]; then echo wow ; fi


