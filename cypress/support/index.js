import "./commands";
import { getMappings } from "../mocks/wiremock";

afterEach(() => {
  cy.resetMocks();
});

Cypress.on("fail", (error) => {
  // On test failure, log all URLs currently mocked
  getMappings()
    .then((response) => response.json())
    .then((json) => {
      json.mappings.forEach((mapping) => {
        Cypress.log({
          name: "mocked URLs",
          displayName: `${mapping.request.method} ${mapping.request.url}`,
          message: mapping.response.body,
        });
      });
    });

  throw error;
});
