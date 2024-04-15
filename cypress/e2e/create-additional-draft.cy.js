describe("Create Additional Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });

    cy.addMock("/lpa-api/v1/persons/130", "GET", {
      status: 200,
      body: {
        id: 130,
        dob: "05/05/1970",
        firstname: "John",
        surname: "Doe",
        uId: "7000-0000-0007",
        addressLine1: "Road",
        town: "City",
        postcode: "A12 3CD",
        country: "GB",
        cases: [
          {
            caseSubtype: "personal-welfare",
            caseType: "DIGITAL_LPA",
            id: 56,
            status: "Draft",
            uId: "M-ABCD-0000-1234",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/donors/130/digital-lpas", "POST", {
      status: 200,
      body: {
        types: ["personal-welfare"],
        source: "PHONE",
      },
    });

    cy.visit("/create-additional-draft-lpa?id=130");
  });

  it("creates an additional digital LPA", () => {
    cy.contains("Create a draft LPA for John Doe");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("The donor, using the address above").click();

    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-panel").contains(
      "2 draft LPAs for John Doe have been created.",
    );
  });
});
