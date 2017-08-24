FROM alpine:3.6
RUN apk --update add tzdata
COPY ./swagchat-api /usr/bin/swagchat-api
EXPOSE 9000
CMD ["/usr/bin/swagchat-api"]
