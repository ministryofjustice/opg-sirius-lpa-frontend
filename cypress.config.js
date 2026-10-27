const { defineConfig } = require("cypress");

module.exports = defineConfig({
  expose: {
    MOCK_SERVER_URI:
      process.env.CYPRESS_MOCK_SERVER_URI || "http://localhost:8563",
  },
  e2e: {
    baseUrl: "http://localhost:8888/",
    supportFile: "cypress/support/index.js",
  },
});
