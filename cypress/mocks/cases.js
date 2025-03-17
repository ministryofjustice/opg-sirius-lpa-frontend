import { addMock } from "./wiremock";

const warnings = {
  async empty(caseId) {
    await addMock(
      `/lpa-api/v1/cases/${caseId}/warnings`,
      "GET",
      {
        status: 200,
        body: [],
      },
      1,
    );
  },
};

const tasks = {
  async empty(caseId) {
    await addMock(
      `/lpa-api/v1/cases/${caseId}/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC`,
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );
  },
};

export { warnings, tasks };
