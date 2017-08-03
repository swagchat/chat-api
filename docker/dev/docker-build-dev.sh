#!/bin/bash

user=betchi
tag=latest
if [ "$1" != "" ]; then
	tag=$1
fi

echo -e "\033[36m----------> Building docker image [$user/alpine-gobuild]\033[0m"
docker build -t $user/alpine-gobuild -f ./Dockerfile-GoBuild .
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit
fi

echo -e "\033[36m----------> Building go binary for alpine linux [swagchat-api]\033[0m"
docker run -i -v $GOPATH/src/github.com/fairway-corp/swagchat-api:/go/src/github.com/fairway-corp/swagchat-api -w /go/src/github.com/fairway-corp/swagchat-api $user/alpine-gobuild go build
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit
fi

mv $GOPATH/src/github.com/fairway-corp/swagchat-api/swagchat-api swagchat-api

echo -e "\033[36m----------> Building docker image [$user/swagchat-api:$tag]\033[0m"
docker build -t $user/swagchat-api:$tag -f ./Dockerfile-Dev .
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit
fi

rm swagchat-api

echo -e "\033[36m----------> Pushing [$user/swagchat-api:$tag]\033[0m"
docker push $user/swagchat-api:$tag
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit
fi

echo -e "\033[36m----------> Complete!\033[0m"

