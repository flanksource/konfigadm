#!/bin/bash
GITHUB_USER=$(echo $GITHUB_REPOSITORY | cut -d/ -f1)
NAME=$(echo $GITHUB_REPOSITORY | cut -d/ -f2)
TAG=$(echo $GITHUB_REF | sed 's|refs/tags/||')


VERSION="v$TAG built $(date)"

make linux darwin windows compress

GO111MODULE=off go get github.com/aktau/github-release
go get github.com/aktau/github-release


if [[ "$SNAPSHOT" == "true" ]]; then
  echo Releasing pre-release
  github-release release -u $GITHUB_USER -r ${NAME} --tag $TAG --pre-release
else
  echo Releasing final release
  github-release release -u $GITHUB_USER -r ${NAME} --tag $TAG
fi

echo Uploading $NAME
github-release upload -R -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME} -f .bin/${NAME}
echo Uploading ${NAME}_osx
github-release upload -R -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}_osx -f .bin/${NAME}_osx
echo Uploading ${NAME}.exe
github-release upload -R -u $GITHUB_USER -r ${NAME} --tag $TAG -n ${NAME}_osx -f .bin/${NAME}.exe
