FROM node:24.22.1-alpine3.23@sha256:0468540aa58203ff9e40c2aab3907c52644ab818af6176b1c70a026137dc92a2

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
