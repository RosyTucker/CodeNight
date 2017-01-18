set -e

export GOBIN=$(pwd)/bin

go get
go install

./bin/codenight.exe