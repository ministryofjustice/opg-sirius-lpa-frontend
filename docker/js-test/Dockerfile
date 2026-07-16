FROM node:24.16.0-alpine3.23@sha256:2bdb65ed1dab192432bc31c95f94155ca5ad7fc1392fb7eb7526ab682fa5bf14

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