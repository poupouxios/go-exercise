# Golang Exercise

This is just an example of how to run everything inside docker for go without the need of having golang install on your machine
Also it will be used to learn Golang.

## Installation

1. First make sure to fetch the latest golang docker with command `docker pull golang:alpine`
2. After that run the image in order to execute commands inside. Also mount the `code` folder to output the files there also. The command is `docker run -it --rm -v ./code:/application -w /application golang-proto /bin/sh`
3. Execute inside the `go mod init golang-exercise` 
4. Now you are ready to create your main.go file and write the code inside.
5. For proto use this command inside the golang-proto docker container to compile the proto `protoc --proto_path=./proto --go_out=paths=source_relative:./exchangeratepb --go-grpc_out=paths=source_relative:./exchangeratepb proto/exchangerate.proto` 