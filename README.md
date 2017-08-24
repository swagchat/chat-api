[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/fairway-corp/swagchat-api)
[![CircleCI](https://circleci.com/gh/fairway-corp/swagchat-chat-api/tree/master.svg?style=svg)](https://circleci.com/gh/fairway-corp/swagchat-chat-api/tree/master)
[![Issue Count](https://lima.codeclimate.com/github/fairway-corp/swagchat-api/badges/issue_count.svg)](https://lima.codeclimate.com/github/fairway-corp/swagchat-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/fairway-corp/swagchat-api)](https://goreportcard.com/report/github.com/fairway-corp/swagchat-api)



# swaghat Chat API

swagchat is an open source chat components for your webapps.

* **Easy to deploy**
* **Easy to customize**
* **Easy to scale**

## Components

* **Chat API (Go)**
* [RTM API (Real Time Messaging) (Go)](https://github.com/fairway-corp/swagchat-rtm-api)
* [Client SDK (TypeScript & JavaScript)](https://github.com/fairway-corp/swagchat-sdk)
* [UIKit (Typescript - React)](https://github.com/fairway-corp/react-swagchat)

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

## Quick start

### Just run the executable binary

You can download binary from [Release page](https://github.com/fairway-corp/swagchat-api/releases)

```
# In the case of macOS
./swagchat-api_darwin_amd64
```

### docker

```
docker pull swagchat/chat-api
docker run swagchat/chat-api
```
[Docker repository](https://hub.docker.com/r/swagchat/chat-api/)

### heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)


## Configuration

To override the default configuration options, make a copy of `config/example.swagchat.yaml` and name it `config/swagchat.yaml`.

Or you can overwrite it with environment variable, runtime parameter.

You can check the variables that can be set with the help command of the executable binary.


```
# In the case of macOS
./swagchat-api_darwin_amd64 -h
```

## Development

### go version

1.7 or higher

### go testing

Only http client test, and not completed yet. Test run with datastore is sqlite3 and storage is local.

```
go test $(go list ./... | grep -v vendor)
```

## Profiling

To display the profile by http request, please set as follows in the setting file.

This is using pprof serves provided by golang.

```
profiling: on
```

### Pprof api list

```
/debug/pprof               pprof portal
/debug/pprof/profile       CPU profile
/debug/pprof/goroutine     goroutine profile
/debug/pprof/heap          heap profile
/debug/pprof/block         blocking profile
/debug/pprof/threadcreate  OS thread profile
```


## License

MIT License.
