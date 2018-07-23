FROM golang:1.10-alpine AS build
LABEL maintainer betchi

RUN apk add --update --no-cache alpine-sdk bash python
WORKDIR /root
RUN git clone -b v0.11.4 https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure
RUN make
RUN make install

COPY . /go/src/github.com/fairway-corp/swag-operator-api
WORKDIR /go/src/github.com/fairway-corp/swag-operator-api/
RUN go build -o swag-operator-api

FROM alpine:3.7
LABEL maintainer betchi

RUN apk --no-cache --update upgrade \
  && apk add --update --no-cache tzdata ca-certificates \
  && update-ca-certificates --fresh

RUN mkdir -p /app
COPY --from=build /go/src/github.com/fairway-corp/swag-operator-api/swag-operator-api /app/swag-operator-api
COPY --from=build /usr/local/lib/librdkafka.a /usr/local/lib/librdkafka.a
COPY --from=build /usr/local/lib/librdkafka.so /usr/local/lib/librdkafka.so
COPY --from=build /usr/local/lib/librdkafka.so.1 /usr/local/lib/librdkafka.so.1
COPY --from=build /usr/local/include/librdkafka /usr/local/include/librdkafka

STOPSIGNAL SIGTERM

EXPOSE 8106 9106
CMD /bin/sh -c "cd /app && ./swag-operator-api"
