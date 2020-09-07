NAME:=konfigadm
NETLIFY_ID:=822e0cb8-b16a-430f-a25c-bebeec0c33d0

ifeq ($(VERSION),)
VERSION=$(shell git describe --tags  --long)$(shell date +"%H%M%S")
endif

all: test integration

.PHONY: clean
clean:
	rm *.img *.vmx *.vmdk *.qcow2 *.ova || true

.PHONY: deps
deps:
	command -v go2xunit 2>&1 > /dev/null || go get github.com/tebeka/go2xunit
	command -v esc 2>&1 > /dev/null || go get -u github.com/mjibson/esc

.PHONY: linux
linux: pack
	GOOS=linux GOARCH=386 go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: darwin
darwin: pack
	GOOS=darwin go build -o ./.bin/$(NAME)_osx -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: windows
windows: pack
	GOOS=windows go build -o ./.bin/$(NAME).exe -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: compress
compress:
	# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
	command -v upx 2>&1 >  /dev/null  || (sudo apt-get update && sudo apt-get install -y xz-utils && wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; tar xf upx.tar.xz; mv upx-3.96-amd64_linux/upx /usr/bin )
	upx ./.bin/$(NAME) ./.bin/$(NAME)_osx ./.bin/$(NAME).exe

.PHONY: install
install:
	go build -ldflags '-X main.version=$(VERSION)' -o konfigadm
	mv konfigadm /usr/local/bin/konfigadm

.PHONY: test
test: deps pack
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
e2e-all: deps linux debian ubuntu20 ubuntu18 ubuntu16 fedora centos

.PHONY: amazonlinux2
amazonlinux2: deps
	IMAGE=quay.io/footloose/amazonlinux2:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: debian9
debian9: deps
	IMAGE=jrei/systemd-debian:9 ./scripts/e2e.sh $(test)

.PHONY: debian10
debian10: deps
	IMAGE=quay.io/footloose/debian10:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: ubuntu16
ubuntu16: deps
	IMAGE=quay.io/footloose/ubuntu16.04:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: ubuntu18
ubuntu18: deps
	IMAGE=quay.io/footloose/ubuntu18.04:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: ubuntu20
ubuntu20: deps
	IMAGE=docker.io/flanksource/ubuntu:20.04 ./scripts/e2e.sh $(test)

.PHONY: fedora29
fedora29: deps
	IMAGE=quay.io/footloose/fedora29:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: photon3
photon3: deps
	IMAGE=docker.io/tarun18/photon:3.0 ./scripts/e2e.sh $(test)

.PHONY: centos7
centos7: deps
	IMAGE=quay.io/footloose/centos7:0.6.3 ./scripts/e2e.sh $(test)

.PHONY: centos8
centos8: deps
	IMAGE=quay.io/footloose/centos8:latest ./scripts/e2e.sh $(test)

.PHONY: docs
docs:
	command -v mkdocs 2>&1 > /dev/null || pip install mkdocs mkdocs-material
	mkdocs build -d build/docs

.PHONY: deploy-docs
deploy-docs: docs
	command -v netlify 2>&1 > /dev/null || sudo npm install -g netlify-cli
	netlify deploy --site $(NETLIFY_ID) --prod --dir build/docs

.PHONY: pack
pack:
	esc --prefix resources/ --ignore "static.go"  -o resources/static.go --pkg resources resources

.PHONY: test-env
test-env:
	docker run --privileged -v /sys/fs/cgroup:/sys/fs/cgroup -v $(PWD):$(PWD) -w $(PWD)  --rm -it quay.io/footloose/debian10:0.6.3 /lib/systemd/systemd
