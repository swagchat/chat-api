#!/bin/bash

rm -f ./protobuf/*.gw.go
rm -f ./protobuf/*.pb.go

protoc \
  -I$GOPATH/src/github.com/swagchat/protobuf \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:./protobuf \
  --grpc-gateway_out=logtostderr=true:./protobuf \
  $GOPATH/src/github.com/swagchat/protobuf/chatMessage.proto \
  $GOPATH/src/github.com/swagchat/protobuf/chatService.proto

mv ./protobuf/github.com/swagchat/protobuf/* ./protobuf/
rm -rf ./protobuf/github.com

protoc-go-inject-tag -input=./protobuf/chatMessage.pb.go

