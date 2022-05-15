rm go-docker

# GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o /Users/project/tools/tools/bin/swag ./cmd/swag/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o go-docker .

upx go-docker
