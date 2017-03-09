VERSION := "1.1.0"
REPO := envsec
USER := kreuzwerker
FLAGS := "-X=main.build=`git rev-parse --short HEAD` -X=main.version=$(VERSION)"

GITHUB_TOKEN ?= $(shell cat .token)
export GITHUB_TOKEN

.PHONY: build clean release retract

build:
	cd bin && gox -osarch="linux/amd64 linux/arm darwin/amd64" -ldflags $(FLAGS) -output "../build/{{.OS}}-{{.Arch}}/es";

clean:
	rm -rf build

release: clean build
	git tag $(VERSION) -f && git push --tags -f
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) --name es-linux --file build/linux-amd64/es
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) --name es-linux-arm --file build/linux-arm/es
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) --name es-osx --file build/darwin-amd64/es

retract:
	github-release delete --tag $(VERSION)
