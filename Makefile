VERSION := 1.0.0
FLAGS := "-X main.build `git rev-parse --short HEAD` -X main.version $(VERSION)"

#export GO15VENDOREXPERIMENT=1
export GITHUB_REPO := envsec
export GITHUB_USER := kreuzwerker
export TOKEN = `cat .token`

.PHONY: build

test:
	go test

build:
	cd bin && \
		gox -os="linux darwin freebsd" -arch="amd64" -ldflags $(FLAGS) -output "../build/{{.OS}}_{{.Arch}}/es";

clean:
	rm -rf build manifest

install:
	go get -v ./...

# generate keypair for signing
keys/secret:
	mkdir -p keys
	signify -G -n -p keys/public -s keys/secret

# generate a SHA2 manifest and sign it
manifest:
	find build -type f | sed "s|^\build/||" | parallel shasum -ba 256 build/{} | sort -t " " -k 2 | signify -Seq -s keys/secret -x manifest -m -

# verify a manifest signature and then the manifest itself
verify:
	signify -Veq -p keys/public -x manifest -m - | shasum -ba 256 -c

.PHONY: install verify

release: clean build manifest
	git tag $(VERSION) -f && git push --tags -f
	github-release release --tag $(VERSION) -s $(TOKEN)
	github-release upload --tag $(VERSION) -s $(TOKEN) --name es-osx --file build/darwin_amd64/es
	github-release upload --tag $(VERSION) -s $(TOKEN) --name es-freebsd --file build/freebsd_amd64/es
	github-release upload --tag $(VERSION) -s $(TOKEN) --name es-linux --file build/build/linux_amd64/es
	github-release upload --tag $(VERSION) -s $(TOKEN) --name manifest --file manifest

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)
