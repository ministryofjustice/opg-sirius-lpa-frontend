FROM node:26.13.0-alpine3.23@sha256:44a3247313f9bc069a4ec25d01b66254ed7dc4d0408727d8720dc827f0753423

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