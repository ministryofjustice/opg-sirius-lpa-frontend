FROM node:alpine3.18

RUN apk add --no-cache git

RUN addgroup -S node && \
    adduser -S -g node node

USER node
