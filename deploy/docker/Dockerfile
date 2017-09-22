FROM alpine:3.6
RUN apk --update add tzdata
COPY ./chat-api /usr/bin/swagchat-chat-api
EXPOSE 9000
CMD ["/usr/bin/swagchat-chat-api"]
