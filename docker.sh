#!/bin/sh

docker build . -t docker.swagchat.io:30000/chat-api && docker push docker.swagchat.io:30000/chat-api