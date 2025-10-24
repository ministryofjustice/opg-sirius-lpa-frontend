import { addMock } from "./wiremock";

async function get(id, body, priority = 1) {
  await addMock(
    `/lpa-api/v1/cases/${id}`,
    "GET",
    {
      status: 200,
      body: body,
    },
    priority,
  );
}

const warnings = {
  async empty(caseId, priority = 1) {
    await addMock(
      `/lpa-api/v1/cases/${caseId}/warnings`,
      "GET",
      {
        status: 200,
        body: [],
      },
      priority,
    );
  },
};

const tasks = {
  async empty(caseId, priority = 1) {
    await addMock(
      `/lpa-api/v1/cases/${caseId}/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC`,
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
      priority,
    );
  },
};

export { get, warnings, tasks };
