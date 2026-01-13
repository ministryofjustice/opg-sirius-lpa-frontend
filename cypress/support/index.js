import "./commands";
import { mappings } from "../mocks/wiremock"

afterEach(() => {
  cy.resetMocks();
});

Cypress.on('fail', (error, runnable) => {
  debugger

  let wiremockMappings = mappings();

  wiremockMappings.then((response) => response.json())
    .then((json) => {
      let output = "Wiremock mappings: \n\n"

      json.mappings.forEach((mapping) => {
        output += mapping.request.method + " " + mapping.request.url + " " + mapping.response.status + "\n"
        output += mapping.response.body + "\n\n"
      })

      console.log(output)
    })

  throw error
})