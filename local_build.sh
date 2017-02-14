set -e

source env.sh

go install

mongod & go run main.go