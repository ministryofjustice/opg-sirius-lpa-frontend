import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  addMock(url, method, response);
});

Cypress.Commands.add("resetMocks", async () => {
  reset();
});
