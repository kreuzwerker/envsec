VERSION := 0.1.0
FLAGS := "-X main.build `git rev-parse --short HEAD` -X main.version $(VERSION)"

export GO15VENDOREXPERIMENT=1

build:
	pushd bin && \
		gox -os="linux darwin windows freebsd" -arch="amd64" -ldflags $(FLAGS) -output "../build/{{.OS}}_{{.Arch}}/es-{{.Dir}}"; \
		popd

clean:
	rm -rf build manifest

install:
	go get -u github.com/mitchellh/gox

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
