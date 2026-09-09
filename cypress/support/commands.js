import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  await addMock(url, method, response);
});

Cypress.Commands.add("resetMocks", async () => {
  await reset();
});

Cypress.Commands.add("mockDocumentFile", (uuid) => {
  cy.intercept("GET", `/lpa-api/v1/documents/${uuid}/file`, {
    statusCode: 200,
    headers: {
      "Content-Type": "application/pdf",
    },
    fixture: "document.pdf",
  }).as(`documentFile_${uuid}`);
});
