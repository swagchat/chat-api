#!/bin/bash

rm -f ./protogen/*.gw.go
rm -f ./protogen/*.pb.go

protoc \
  -I$GOPATH/src/github.com/swagchat/chat-api/proto \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:./protogen \
  --grpc-gateway_out=logtostderr=true:./protogen \
  $GOPATH/src/github.com/swagchat/chat-api/proto/chatMessage.proto \
  $GOPATH/src/github.com/swagchat/chat-api/proto/chatService.proto

