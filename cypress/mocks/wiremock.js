async function addMock(url, method, response, priority) {
  if (typeof response.body !== "string") {
    response.body = JSON.stringify(response.body);
  }

  await fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings`, {
    method: "POST",
    body: JSON.stringify({
      request: {
        url,
        method,
      },
      response,
      priority,
    }),
  });
}

async function reset() {
  await fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings/reset`, {
    method: "POST",
  });
}

export { addMock, reset };
