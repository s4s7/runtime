#!/bin/bash
set -e -o pipefail

go test -race -timeout 20m -v ./...

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: ${GOCOVMODE-atomic}" > coverage.txt

# Standard go tooling behavior is to ignore dirs with leading underscores
# skip generator for race detection and coverage
for dir in $(go list ./...)
do
  pth="$GOPATH/src/$dir"
  go test -covermode=${GOCOVMODE-atomic} -coverprofile=${pth}/profile.out $dir
  if [ -f $pth/profile.out ]
  then
      cat $pth/profile.out | tail -n +2 >> coverage.txt
      rm $pth/profile.out
  fi
done

go tool cover -func coverage.txt
