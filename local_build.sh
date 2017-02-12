set -e

source env.sh
./get_dependencies.sh

go install

mongod & go run main.go