[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/swagchat/chat-api)](https://goreportcard.com/report/github.com/swagchat/chat-api)
[![CircleCI](https://circleci.com/gh/swagchat/chat-api/tree/master.svg?style=svg)](https://circleci.com/gh/swagchat/chat-api/tree/master)
[![Maintainability](https://api.codeclimate.com/v1/badges/5c3261e99582f147950c/maintainability)](https://codeclimate.com/github/swagchat/chat-api/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/5c3261e99582f147950c/test_coverage)](https://codeclimate.com/github/swagchat/chat-api/test_coverage)

# swaghat Chat API

swagchat is an open source chat components for your webapps.

chat-api is designed to be easy to introduce to your microservices as well.

**Currently developing for version 1**

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fswagchat%2Fchat-api.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fswagchat%2Fchat-api?ref=badge_shield)


##### Repository structure

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

The agent has various configuration options that can be specified via the command-line or via configuration files. All of the configuration options are completely optional. Defaults are specified with their descriptions.

Configuration precedence is evaluated in the following order:

1. Command line arguments
1. Environment Variables
1. Configuration files

### Specify the setting file (yaml format)

To override the default configuration options, make a copy of `defaultConfig.yaml` and then specify that file name in runtime parameter `config` and execute.

```
./chat-api -config myConfig.yaml
```

### Specify environment variables

You can overwrite it with environment variable.

```
export HTTP_PORT=80 && ./chat-api
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


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fswagchat%2Fchat-api.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fswagchat%2Fchat-api?ref=badge_large)