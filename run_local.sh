set -e

export GOBIN=$(pwd)/bin

go get
go install


source env.sh
./bin/codenight.exe