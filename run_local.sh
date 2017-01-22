set -e

export GOBIN=$(pwd)/bin

go get
go install


source env.sh
mongod & ./bin/codenight.exe