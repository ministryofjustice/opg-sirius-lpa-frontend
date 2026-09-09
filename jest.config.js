module.exports = {
  testEnvironment: "jsdom",
  transform: {
    "^.+\\.js$": "babel-jest",
  },
  testMatch: ["**/*.test.js"],
  moduleFileExtensions: ["js"],
  transformIgnorePatterns: ["node_modules/(?!(pdfjs-dist)/)"],
  setupFilesAfterEnv: [],
};
