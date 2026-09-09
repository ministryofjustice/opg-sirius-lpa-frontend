FROM node:24.20.0-alpine3.23@sha256:0388af2af070cd4736a1567cfed02469ba117848845b4165d87a333edb53d2ca

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
