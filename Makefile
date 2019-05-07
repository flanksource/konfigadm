all: test integration

.PHONY: test
test:
	go test -v ./pkg/... ./cmd/... -race -coverprofile=coverage.txt -covermode=atomic


.PHONY: integration
integration:
	go test -v ./test -race -coverprofile=coverage.txt -covermode=atomic



