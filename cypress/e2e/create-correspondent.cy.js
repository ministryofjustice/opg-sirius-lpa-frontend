describe("Create correspondent", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/create-correspondent?id=1&caseId=2");
  });

  it("creates a correspondent on an EPA", () => {
    cy.contains("Add a correspondent");
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
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-epa");
  });
});
