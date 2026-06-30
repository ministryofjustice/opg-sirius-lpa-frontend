FROM node:24.17.0-alpine3.23@sha256:7c70d1235c0b4c2bc9eeed5393d19f1bbdde6885ba0d58ba62bb385d7b0f3ff1

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

RUN npm i --ignore-scripts -g pa11y-ci@2.4 @lhci/cli@0.4.x

COPY pa11yci.json .pa11yci
COPY lighthouserc.js .

CMD ["pa11y-ci"]
