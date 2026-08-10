FROM node:24.18.1-alpine3.23@sha256:c2cc26d8f991c2db236ad51a61efee843c482372d6d22570787309d511694110

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