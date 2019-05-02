build:
	go build .

release-test:
	goreleaser --rm-dist --snapshot --skip-publish --skip-validate release
