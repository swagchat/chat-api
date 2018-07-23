#!/bin/sh

docker build . -t docker.swagchat.io:30000/swag-operator-api && docker push docker.swagchat.io:30000/swag-operator-api