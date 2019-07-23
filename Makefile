all: test docs integration

.PHONY: clean
clean:
	rm *.img *.vmx *.vmdk *.qcow2 *.ova || true

.PHONY: deps
deps:
	GO111MODULE=off which go2xunit 2>&1 > /dev/null || go get github.com/tebeka/go2xunit

.PHONY: linux
linux:
	 CGO_ENABLED=0 GOOS=linux go build -o dist/konfigadm.tmp -ldflags '-X main.version=built-$(shell date +%Y%m%d%M%H%M%S)' . && mv dist/konfigadm.tmp dist/konfigadm

.PHONY: test
test: deps
	mkdir -p test-output
	go test -v ./pkg/... ./cmd/... -race -coverprofile=coverage.txt -covermode=atomic | tee unit.out
	cat unit.out | go2xunit --fail -output test-output/unit.xml
	rm unit.out

.PHONY: integration
integration: linux
		./scripts/e2e.sh $(test)

.PHONY: e2e
e2e: linux
		./scripts/e2e.sh $(test)

.PHONY: e2e-all
e2e-all: deps linux debian ubuntu ubuntu16 fedora centos

.PHONY: debian9
debian9: deps
		IMAGE=jrei/systemd-debian:9 ./scripts/e2e.sh $(test)

.PHONY: debian
debian: deps
		IMAGE=jrei/systemd-debian:latest ./scripts/e2e.sh $(test)

.PHONY: ubuntu16
ubuntu16: deps
		IMAGE=jrei/systemd-ubuntu:16.04 ./scripts/e2e.sh $(test)

.PHONY: ubuntu
ubuntu: deps
		IMAGE=jrei/systemd-ubuntu:18.04 ./scripts/e2e.sh $(test)

.PHONY: fedora
fedora: deps
		IMAGE=jrei/systemd-fedora:latest ./scripts/e2e.sh $(test)

.PHONY: centos
centos: deps
		IMAGE=jrei/systemd-centos:7 ./scripts/e2e.sh $(test)

.PHONY: docs
docs:
	git remote add docs "https://$(GH_TOKEN)@github.com/moshloop/konfigadm.git"
	git fetch docs && git fetch docs gh-pages:gh-pages
	mkdocs gh-deploy -v --remote-name docs -m "Deployed {sha} with MkDocs version: {version} [ci skip]"
