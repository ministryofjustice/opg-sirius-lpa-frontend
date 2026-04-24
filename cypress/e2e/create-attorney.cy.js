describe("Create Attorney", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/reference-data/relationshipToDonor", "GET", {
      status: 200,
      body: [
        {
          handle: "NO RELATION",
          label: "No relation",
        },
        {
          handle: "OTHER RELATION",
          label: "Other relation",
        },
      ],
    });

    cy.addMock("/lpa-api/v1/epas/2/attorneys", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/create-attorney?id=1&caseId=2");

    cy.get("#f-salutation").type("Prof");
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").type("Josefina");
    cy.get("#f-surname").type("Vanvolkenburg");
    cy.get("#f-dob").type("1978-04-19");
    cy.get("#f-companyName").type("ACME");
    cy.get('[data-module="app-address-finder"] .govuk-link').click();
    cy.get("#f-addressLine1").type("29737 Andrew Plaza");
    cy.get("#f-addressLine2").type("Apt. 814");
    cy.get("#f-addressLine3").type("Gislasonside");
    cy.get("#f-town").type("Hirthehaven");
    cy.get("#f-county").type("Saskatchewan");
    cy.get("#f-postcode").type("S7R 9F9");
    cy.get("#f-country").type("Canada");
    cy.get("label[for=f-isAirmailRequired]").click();
    cy.get("#f-phoneNumber").type("072345678");
    cy.get("#f-email").type("m.vancolkenburg@ca.test");
    cy.get("#f-relationshipToDonor").select("Other relation");
    cy.get("label[for=f-isAttorneyActive]").click();
  });

  it("creates an attorney on an EPA", () => {
    cy.contains("Add an attorney");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-epa");
  });

  it("creates an attorney on an EPA and add another attorney", () => {
    cy.contains("Add an attorney");
    cy.get("input[type=submit][name=add-another]").click();
    cy.url().should("include", "create-attorney");
  });
});
