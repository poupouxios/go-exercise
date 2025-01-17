# Use an official Golang image as the base
FROM golang:1.23.4

# Install the required tools for protoc
RUN apt-get update && apt-get install -y \
    unzip \
    curl

# Install protoc
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-x86_64.zip \
    && unzip protoc-23.4-linux-x86_64.zip -d /usr/local \
    && rm -f protoc-23.4-linux-x86_64.zip

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
