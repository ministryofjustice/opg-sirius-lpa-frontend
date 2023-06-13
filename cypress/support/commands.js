Cypress.Commands.add("addMock", async (url, method, response) => {
  if (typeof response.body !== "string") {
    response.body = JSON.stringify(response.body);
  }

  await fetch("http://localhost:1234/__admin/mappings", {
    method: "POST",
    body: JSON.stringify({
      request: {
        url,
        method,
      },
      response,
    }),
  });
});

Cypress.Commands.add("resetMocks", async () => {
  await fetch("http://localhost:1234/__admin/mappings/reset", {
    method: "POST",
  });
});
