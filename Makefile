all: test integration


.PHONY: linux
linux:
	GOOS=linux go build -o dist/configadm -ldflags '-X main.version=built-$(shell date +%Y%m%d%M%H%M%S)' .

.PHONY: test
test:
	go test -v ./pkg/... ./cmd/... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: integration
integration: linux
	go test -v ./test -race -coverprofile=integ.txt -covermode=atomic



