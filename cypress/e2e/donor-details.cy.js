describe("Edit dates", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {
        id: 1,
        uId: "7000-0000-0013",
        email: "consuela@somesite.example",
        dob: "25/03/1960",
        salutation: "Dr",
        addressLine1: "3 Church Road",
        addressLine2: "",
        addressLine3: "",
        town: "Blackpool",
        county: "Lancashire",
        postcode: "FY48 7CY",
        country: "United Kingdom",
        sageId: "L0000001",
      },
    });
    cy.visit("/donor/1/details");
  });

  it("edits the dates", () => {
    cy.get('[data-cy="donorSalutation"]').contains("Dr");
    cy.get('[data-cy="donorDOB"]').contains("25/03/1960");
    cy.get('[data-cy="donorRecord"]').contains("7000-0000-0013");
    cy.get('[data-cy="donorEmail"]').contains("consuela@somesite.example");
    cy.get('[data-cy="donorAddress"]').contains(
      "3 Church Road, Blackpool, Lancashire, United Kingdom, FY48 7CY",
    );
  });
});
