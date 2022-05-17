rm go-docker

# GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o /Users/project/tools/tools/bin/swag ./cmd/swag/main.go
#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o go-docker .

#CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CGO_LDFLAGS="-static"  CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++  go build -o go-docker .

upx go-docker
