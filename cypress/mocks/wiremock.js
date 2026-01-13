async function addMock(url, method, response, priority = 1) {
  let presignMock = null;

  if (
    method === "GET" &&
    url.match(/^\/lpa-api\/v1\/digital-lpas\/M(-[A-Z0-9]{4}){3}\/?$/)
  ) {
    // if we need to mock this route there is a good chance the test hits
    // /lpa-details, and will therefore need to also mock with query
    // ?presignImages, but assign a lower priority, so it can be overwritten
    presignMock = addMock(url + "?presignImages", "GET", response, 2);
  }

  if (typeof response.body !== "string") {
    response.body = JSON.stringify(response.body);
  }

  const mock = fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings`, {
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

  await Promise.allSettled([presignMock, mock]);

  return mock;
}

async function reset() {
  await fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings/reset`, {
    method: "POST",
  });
}

function mappings() {
  return fetch(`${Cypress.env("MOCK_SERVER_URI")}/__admin/mappings`);
}

export { addMock, reset, mappings };
