#!/bin/bash

rm -f ./*.gw.go
rm -f ./*.pb.go

protoc \
  -I./ \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:../../../ \
  --grpc-gateway_out=logtostderr=true:../../../ \
  chatMessage.proto \
  chatService.proto

