FROM node:24.21.0-alpine3.23@sha256:159fe64649038c30f8cc1ec4be3af3a6e93e3648678c31294e2c5058dbeb99f3

RUN apk add --no-cache \
    ca-certificates \
    chromium \
    harfbuzz \
    freetype \
    freetype-dev \
    nodejs \
    nss \
    ttf-freefont

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts

COPY pa11yci.json .pa11yci
