import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  // if we need to mock this route there is a good chance the test hits
  // /lpa-details, and will therefore need to also mock with query
  // ?presignImages, but assign a lower priority so it can be overwritten
  if (
    method == "GET" &&
    url.match(/^\/lpa-api\/v1\/digital-lpas\/M(-[A-Z0-9]{4}){3}$/)
  ) {
    await addMock(url + "?presignImages", method, response, 2);
  }

  await addMock(url, method, response, 1);
});

Cypress.Commands.add("resetMocks", async () => {
  await reset();
});
