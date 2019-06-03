all: test docs integration


.PHONY: linux
linux:
	GOOS=linux go build -o dist/konfigadm -ldflags '-X main.version=built-$(shell date +%Y%m%d%M%H%M%S)' .

.PHONY: test
test:
	mkdir -p test-output
	go test -v ./pkg/... ./cmd/... -race -coverprofile=coverage.txt -covermode=atomic | tee unit.out
	cat unit.out | go2xunit -output test-output/unit.xml
	rm unit.out

.PHONY: integration
integration: linux
		./scripts/e2e.sh $(test)

.PHONY: e2e
e2e: linux
		./scripts/e2e.sh $(test)

.PHONY: e2e-all
e2e-all: linux debian ubuntu ubuntu16 fedora centos

.PHONY: debian9
debian9:
		IMAGE=jrei/systemd-debian:9 ./scripts/e2e.sh $(test)

.PHONY: debian
debian:
		IMAGE=jrei/systemd-debian:latest ./scripts/e2e.sh $(test)

.PHONY: ubuntu16
ubuntu16:
		IMAGE=jrei/systemd-ubuntu:16.04 ./scripts/e2e.sh $(test)

.PHONY: ubuntu
ubuntu:
		IMAGE=jrei/systemd-ubuntu:18.04 ./scripts/e2e.sh $(test)

.PHONY: fedora
fedora:
		IMAGE=jrei/systemd-fedora:latest ./scripts/e2e.sh $(test)

.PHONY: centos
centos:
		IMAGE=jrei/systemd-centos:7 ./scripts/e2e.sh $(test)

.PHONY: docs
docs:
	git remote add docs "https://$(GH_TOKEN)@github.com/moshloop/konfigadm.git"
	git fetch docs && git fetch docs gh-pages:gh-pages
	mkdocs gh-deploy -v --remote-name docs -m "Deployed {sha} with MkDocs version: {version} [ci skip]"
