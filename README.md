[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/swagchat/chat-api)](https://goreportcard.com/report/github.com/swagchat/chat-api)

# swaghat Chat API

swagchat is an open source chat components for your webapps.

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")


##### Related repositories

* [RTM API (Real Time Messaging API)](https://github.com/swagchat/rtm-api)
* [SDK (TypeScript & JavaScript)](https://github.com/swagchat/swagchat-sdk-js)
* [UIKit (A set of React components)](https://github.com/swagchat/react-swagchat)


## API Reference

Currently writing in OAI 3

## Datastore

You can choose from the followings.

* sqlite3
* MySQL
* Google Cloud SQL

## Storage

You can choose from the followings.

* Local Filesystem
* Google Cloud Storage
* Amazon S3

## Quick start

### Just run the executable binary

You can download binary from [Release page](https://github.com/swagchat/chat-api/releases)

```
# In the case of macOS
./swagchat-api_darwin_amd64
```

### docker

```
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
