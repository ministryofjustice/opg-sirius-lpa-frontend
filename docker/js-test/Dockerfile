FROM node:26.13.1-alpine3.23@sha256:39d8c143e60e99ae98f014bebea1e72fc96d3c9bdf9504bc3a27af5c05bfa39f

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