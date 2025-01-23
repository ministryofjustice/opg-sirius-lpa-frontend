import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  await addMock(url, method, response);
});

Cypress.Commands.add("resetMocks", async () => {
  await reset();
});
