import { addMock } from "./wiremock";

async function empty(caseId) {
  await addMock(`/lpa-api/v1/cases/${caseId}/warnings`, "GET", {
    status: 200,
    body: [],
  });
}

export { empty };
