import { addMock } from "./wiremock";

function empty(caseId) {
  addMock(`/lpa-api/v1/cases/${caseId}/warnings`, "GET", {
    status: 200,
    body: [],
  });
}

export { empty };
