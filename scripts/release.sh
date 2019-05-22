#!/bin/bash
NAME=konfigadm
GITHUB_USER=moshloop
TAG=$(git tag --points-at HEAD )

[[ -e "checkout.txt" ]] && rm checkout.txt
go get github.com/goreleaser/goreleaser
goreleaser release --rm-dist

go get github.com/aktau/github-release
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME} -f dist/linux_amd64/${NAME}
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}_osx -f dist/darwin_amd64/${NAME}
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}.exe -f dist/windows_amd64/${NAME}.exe
