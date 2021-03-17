GOPKG ?=	moul.io/zapring
DOCKER_IMAGE ?=	moul/zapring
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate:
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	go doc -all > .tmp/usage.txt
	embedmd -w README.md
	rm -rf .tmp
.PHONY: generate

lint:
	cd tool/lint; make
.PHONY: lint
