FROM swagchat/build-base:1.0.0 AS build
LABEL maintainer betchi

COPY . /go/src/github.com/swagchat/chat-api
WORKDIR /go/src/github.com/swagchat/chat-api/
RUN go build -o chat-api

FROM swagchat/deploy-base:1.0.0
LABEL maintainer betchi

RUN apk --no-cache --update upgrade \
  && apk add --update --no-cache tzdata ca-certificates \
  && update-ca-certificates --fresh

RUN mkdir -p /app
COPY --from=build /go/src/github.com/swagchat/chat-api/chat-api /app/chat-api

STOPSIGNAL SIGTERM

EXPOSE 8101 9101
ENTRYPOINT ["/app/chat-api"]
