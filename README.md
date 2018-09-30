[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/swagchat/chat-api)](https://goreportcard.com/report/github.com/swagchat/chat-api)
[![CircleCI](https://circleci.com/gh/swagchat/chat-api/tree/master.svg?style=svg)](https://circleci.com/gh/swagchat/chat-api/tree/master)

# swaghat Chat API

swagchat is an open source chat components for your webapps.

chat-api is designed to be easy to introduce to your microservices as well.

**Currently developing for version 1**

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")


##### Related repositories

* [RTM API (Real Time Messaging API)](https://github.com/swagchat/rtm-api)
* [SDK (TypeScript & JavaScript)](https://github.com/swagchat/swagchat-sdk-js)
* [UIKit (A set of React components)](https://github.com/swagchat/react-swagchat)


## API Reference

[swagger (OpenAPI 2.0)](https://app.swaggerhub.com/apis/swagchat/swagchat-res_tful_api/0.3.2)

Sorry, maintenance is not keeping up.


## Multiple datastore

You can choose from the followings.

* sqlite3
* MySQL
* Google Cloud SQL

## Multiple storage

You can choose from the followings.

* Local Filesystem
* Google Cloud Storage
* Amazon S3

## Multiple tracer

You can choose from the followings.

* jaeger
* zipkin
* elastic APM

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

## Configuration

### Specify the setting file (yaml format)

To override the default configuration options, make a copy of `defaultConfig.yaml` and then specify that file name in runtime parameter `config` and execute.

```
./chat-api -config myConfig.yaml
```

### Specify environment variables

You can overwrite it with environment variable.

```
export HTTP_PORT=80 && ./swagchat-api
```

### Specify runtime parameters

You can overwrite it with runtime parameters.

```
./chat-api -httpPort 80
```

You can check the variables that can be set with the help command of the executable binary.

```
./chat-api -h
```

## Development

### go version

1.8 or higher

## License

MIT License.
