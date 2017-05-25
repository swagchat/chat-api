[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/fairway-corp/swagchat-api)
[![CircleCI](https://circleci.com/gh/fairway-corp/swagchat-api.svg?style=shield&circle-token=06b2dbd153b46662683bb01168a3d13891922252)](https://circleci.com/gh/fairway-corp/swagchat-api)
[![Issue Count](https://lima.codeclimate.com/github/fairway-corp/swagchat-api/badges/issue_count.svg)](https://lima.codeclimate.com/github/fairway-corp/swagchat-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/fairway-corp/swagchat-api)](https://goreportcard.com/report/github.com/fairway-corp/swagchat-api)



# SwagChat RESTful API

SwagChat is an open source chat components for your webapps.

* **Easy to deploy**
* **Easy to customize**
* **Easy to scale**

## Components

* **RESTful API Server (Go) ---> This repository**
* Realtime Messaging (Go) ---> In development ...
* [Client SDK (TypeScript & JavaScript)](https://github.com/fairway-corp/swagchat-sdk)
* UIKit (Typescript - React) ---> In development ...


## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/architecture-201703011307.png "Architecture")

## API Reference

### SWAGGER HUB

[https://app.swaggerhub.com/api/fairway-corp/swagchat-api](https://app.swaggerhub.com/api/fairway-corp/swagchat-api)

### Apiary

[http://docs.swagchat.apiary.io](http://docs.swagchat.apiary.io)

## Datastore

You can choose from the followings.

* sqlite3
* MySQL
* Google Cloud SQL
* Oracle (In development ...)

## Storage

You can choose from the followings.

* Local Filesystem
* Google Cloud Storage
* Amazon S3

## Feature

### Chat
* 1-on-1 Chat
* Group Chat
* Display chat room list

### Message Content Type

* text
* image

## Release Binary

You can download from [here](https://github.com/fairway-corp/swagchat-api/releases)

**Supported OS**

* Windows (amd64)
* Linux (amd64)
* macOS (amd64)

## Quick start

Just run the executable binary.

ex) `./swagchat-api_darwin_amd64`

Not For Production Use!

## Configuration

To override the default configuration options, make a copy of `example.swagchat.yaml` and name it `swagchat.yaml`.

Or you can overwrite it with environment variable, runtime parameter.

You can check the variables that can be set with the help command of the executable binary.

ex) `./swagchat-api_darwin_amd64 -h`

## Development

### go version

1.7 or higher

### go testing

Only http client test, and not completed yet. Test run with datastore sqlite3, storage local.

`go test $(go list ./... | grep -v vendor)`

### go profiling

To view all available profiles, open http://localhost:6060/debug/pprof/ in your browser.

## License

MIT License.
