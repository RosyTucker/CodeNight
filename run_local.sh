set -e

export GOBIN=$(pwd)/bin
export GOPATH=$(pwd)

go install

./bin/Codenight.exe