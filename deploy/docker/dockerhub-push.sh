#!/bin/bash

if [ $# != 0 ] && [ $# != 3 ]; then
	echo -e "\033[35mInvalid arguments. ex) ./dockerhub-push.sh swagchat api latest\033[0m"
	exit 1
fi

user=swagchat
image=api
tag=latest
if [ "$1" != "" ]; then
	user=$1
fi
if [ "$2" != "" ]; then
	image=$2
fi
if [ "$3" != "" ]; then
	tag=$3
fi

echo -e "\033[36m----------> Building docker image [$user/alpine-gobuild]\033[0m"
docker build -t $user/alpine-gobuild -f ./Dockerfile-GoBuild .
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit 1
fi

echo -e "\033[36m----------> Building go binary for alpine linux [swagchat-api]\033[0m"
docker run -i -v $GOPATH/src/github.com/fairway-corp/swagchat-api:/go/src/github.com/fairway-corp/swagchat-api -w /go/src/github.com/fairway-corp/swagchat-api $user/alpine-gobuild go build
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit 1
fi

mv $GOPATH/src/github.com/fairway-corp/swagchat-api/swagchat-api swagchat-api

echo -e "\033[36m----------> Building docker image [$user/$image:$tag]\033[0m"
docker build -t $user/$image:$tag -f ./Dockerfile-Dev .
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit 1
fi

echo -e "\033[36m----------> Pushing [$user/$image:$tag]\033[0m"
docker push $user/$image:$tag
if [ $? -gt 0 ]; then
	echo -e "\033[35mFailed!\033[0m"
	exit 1
fi

mv swagchat-api swagchat-api_alpine_amd64

echo -e "\033[36m----------> Complete!\033[0m"

