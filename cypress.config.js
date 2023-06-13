const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    baseUrl: "http://localhost:8888/",
    supportFile: "cypress/support/index.js",
  },
});
