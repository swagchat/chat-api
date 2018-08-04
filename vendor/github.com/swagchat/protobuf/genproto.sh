#!/bin/bash

rm -f ./protoc-gen-go/*
rm -f ./protoc-gen-js/*pb.js
rm -f ./protoc-gen-js/*pb.d.ts

protoc \
  -I./proto \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/gogo/protobuf \
  --gogo_out=plugins=grpc:../../../ \
  --grpc-gateway_out=logtostderr=true:../../../ \
  --js_out=import_style=closure,library=swagchatpb,binary:./protoc-gen-js/ \
  --ts_out=./protoc-gen-js/ \
  blockUserMessage.proto \
  blockUserService.proto \
  commonMessage.proto \
  deviceMessage.proto \
  deviceService.proto \
  messageMessage.proto \
  messageService.proto \
  roomMessage.proto \
  roomService.proto \
  roomUserMessage.proto \
  roomUserService.proto \
  userMessage.proto \
  userService.proto \
  userRoleMessage.proto \
  userRoleService.proto \
  webhookService.proto
