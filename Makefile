all: test integration

.PHONY: test
test:
	go test -v ./pkg/... ./cmd/... -race -coverprofile=coverage.txt -covermode=atomic


.PHONY: integration
integration:
	GOOS=linux go build -o dist/configadm -ldflags '-X main.version=built' .
	go test -v ./test -race -coverprofile=coverage.txt -covermode=atomic



