TAG:=$(shell git describe --abbrev=0 --tags)

.PHONY: setup build release

setup:
	go get \
		golang.org/x/vgo \
		github.com/Songmu/goxz/cmd/goxz \
		github.com/tcnksm/ghr

build: dist/$(TAG)

dist/$(TAG): vendor
	CGO_ENABLED=0 goxz -d dist/$(TAG) -z -os darwin,linux -arch amd64,386

vendor:
	vgo mod -vendor

release: dist/$(TAG)
	ghr -u dozen -r mackerel-plugin-sqsd $(TAG) dist/$(TAG)
