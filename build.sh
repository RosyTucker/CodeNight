set -e

go get

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o pkg/codenight-amd64 main.go
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o pkg/codenight-linux-386 main.go
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o pkg/codenight-windows-386.exe main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o pkg/codenight-windows-amd64.exe main.go