FROM cypress/included:10.6.0

WORKDIR /root

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://app:8888

COPY cypress.config.js .
COPY cypress cypress
