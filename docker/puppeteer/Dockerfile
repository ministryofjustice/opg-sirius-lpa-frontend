FROM node:22.12.0-alpine3.19

RUN apk add --no-cache \
      ca-certificates \
      chromium \
      harfbuzz \
      freetype \
      freetype-dev \
      nodejs \
      nss \
      ttf-freefont \
      yarn

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

RUN npm i --ignore-scripts -g pa11y-ci@2.4 @lhci/cli@0.4.x

COPY pa11yci.json .pa11yci
COPY lighthouserc.js .

CMD ["pa11y-ci"]
