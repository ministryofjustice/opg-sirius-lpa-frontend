import { addMock } from "./wiremock";

const warnings = {
  async empty(caseId) {
    await addMock(`/lpa-api/v1/cases/${caseId}/warnings`, "GET", {
      status: 200,
      body: [],
    });
  },
};

export { warnings };
