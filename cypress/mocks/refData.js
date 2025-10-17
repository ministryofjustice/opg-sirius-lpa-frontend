import { addMock } from "./wiremock";

const countries = {
  async get(countriesArray, priority = 1) {
    await addMock(
      `/lpa-api/v1/reference-data/country`,
      "GET",
      {
        status: 200,
        body: countriesArray,
      },
      priority,
    );
  },
  async gbOnly(priority = 1) {
    await addMock(
      `/lpa-api/v1/reference-data/country`,
      "GET",
      {
        status: 200,
        body: [{ handle: "GB", label: "Great Britain" }],
      },
      priority,
    );
  },
};

export { countries };
