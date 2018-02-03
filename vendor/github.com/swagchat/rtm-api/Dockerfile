FROM alpine:3.6
ARG API_VERSION="0.3.0"
ARG EXEC_FILE_NAME="swagchat-rtm-api_alpine_amd64"
RUN apk --update add tzdata curl \
  && curl -LJO https://github.com/swagchat/rtm-api/releases/download/v${API_VERSION}/${EXEC_FILE_NAME} \
  && chmod 700 ${EXEC_FILE_NAME} \
  && mv ${EXEC_FILE_NAME} /bin/swagchat-rtm-api
EXPOSE 9000
CMD ["swagchat-rtm-api"]
