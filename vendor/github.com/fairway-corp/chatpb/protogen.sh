#!/bin/bash

rm -f ./*.gw.go
rm -f ./*.pb.go
rm -f ./*.gen.go

protoc \
  -I./proto \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/gogo/protobuf \
  -I$GOPATH/src/github.com/swagchat/protobuf/proto \
  --gogo_out=plugins=grpc:../../../ \
  --grpc-gateway_out=logtostderr=true:../../../ \
  botConnectorMessage.proto \
  botConnectorService.proto \
  indexMessage.proto \
  indexService.proto \
  messageConnectorMessage.proto \
  messageConnectorService.proto \
  operatorMessage.proto \
  operatorService.proto
