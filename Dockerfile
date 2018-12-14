FROM golang:1.11.2-alpine AS build
LABEL maintainer betchi

ENV LIBRDKAFKA_VERSION 0.11.6

RUN apk add --update --no-cache alpine-sdk bash python
WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN git checkout -b v$LIBRDKAFKA_VERSION refs/tags/v$LIBRDKAFKA_VERSION
RUN ./configure && make && make install

WORKDIR /go/src/github.com/swagchat/chat-api/
COPY Gopkg.toml Gopkg.lock ./
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v -vendor-only=true
COPY . .
RUN go build -o chat-api

FROM alpine:3.7
LABEL maintainer betchi

RUN apk --no-cache --update upgrade \
  && apk add --update --no-cache tzdata ca-certificates \
  && update-ca-certificates --fresh

RUN mkdir -p /app
COPY --from=build /go/src/github.com/swagchat/chat-api/chat-api /app/chat-api
COPY --from=build /go/src/github.com/swagchat/chat-api/defaultConfig.yaml /app/defaultConfig.yaml
COPY --from=build /usr/local/lib/librdkafka.a /usr/local/lib/librdkafka.a
COPY --from=build /usr/local/lib/librdkafka.so /usr/local/lib/librdkafka.so
COPY --from=build /usr/local/lib/librdkafka.so.1 /usr/local/lib/librdkafka.so.1
COPY --from=build /usr/local/include/librdkafka /usr/local/include/librdkafka

ENTRYPOINT ["/app/chat-api"]
CMD ["-grpcPort", "9101"]
