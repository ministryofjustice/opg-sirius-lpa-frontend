FROM node:24.18.0-alpine3.23@sha256:595398b0081eacda8e1c4c5b97b76cd1020e4d58a8ebcb4843b9bca1e79e7436

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    app

WORKDIR /app

COPY package.json .
COPY package-lock.json .
COPY jest.config.js .
COPY .babelrc.js .

RUN npm ci --ignore-scripts

COPY web/assets web/assets

ENTRYPOINT ["npm", "run", "jsTest"]