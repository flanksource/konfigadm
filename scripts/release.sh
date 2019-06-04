#!/bin/bash
NAME=$(basename $(git remote get-url origin | sed 's/\.git//'))
GITHUB_USER=$(basename $(dirname $(git remote get-url origin | sed 's/\.git//')))
TAG=$(git tag --points-at HEAD )

GO111MODULE=off go get github.com/goreleaser/goreleaser
git stash
git clean -fd
docker run --rm --privileged -e GITHUB_TOKEN=$GITHUB_TOKEN -v $PWD:$PWD -v /var/run/docker.sock:/var/run/docker.sock -w $PWD goreleaser/goreleaser release --rm-dist

GO111MODULE=off go get github.com/aktau/github-release
go get github.com/aktau/github-release
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME} -f dist/linux_amd64/${NAME}
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}_osx -f dist/darwin_amd64/${NAME}
github-release upload -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}.exe -f dist/windows_amd64/${NAME}.exe
\
