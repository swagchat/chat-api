FROM alpine:3.6
ARG API_VERSION="0.8.6"
ARG EXEC_FILE_NAME="swagchat-api_alpine_amd64"
RUN apk --update add tzdata curl \
  && apk --no-cache add ca-certificates && update-ca-certificates --fresh \
  && curl -LJO https://github.com/fairway-corp/swagchat-api/releases/download/v${API_VERSION}/${EXEC_FILE_NAME} \
  && chmod 700 ${EXEC_FILE_NAME} \
  && mv ${EXEC_FILE_NAME} /bin/swagchat-api
EXPOSE 9000
CMD ["swagchat-api"]
