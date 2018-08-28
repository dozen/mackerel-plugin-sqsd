TAG:=$(shell git describe --abbrev=0 --tags)

.PHONY: setup build release

setup:
	go get \
		golang.org/x/vgo \
		github.com/Songmu/goxz/cmd/goxz \
		github.com/tcnksm/ghr

build: dist/$(TAG)

goget:
	go get -v

dist/$(TAG): goget
	GO111MODULE=on CGO_ENABLED=0 goxz -d dist/$(TAG) -z -os darwin,linux -arch amd64,386

release: dist/$(TAG)
	ghr -u dozen -r mackerel-plugin-sqsd $(TAG) dist/$(TAG)
