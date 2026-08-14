describe("Create trust corporation form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/trust-corporation", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/create-trust-corporation?id=1&caseId=2");
  });

  it("can create a trust corporation", () => {
    cy.contains("Add a trust corporation");
    cy.get("label[for=f-isReplacementAttorney]").click();
    cy.get("label[for=f-isTrustCorporationActive]").click();
    cy.get("#f-companyName").type("ACME");
    cy.get("#f-companyNumber").type("123");
    cy.get(".govuk-details__summary").click();
    cy.get("#f-addressLine1").type("29737 Andrew Plaza");
    cy.get("#f-addressLine2").type("Apt. 814");
    cy.get("#f-addressLine3").type("Gislasonside");
    cy.get("#f-town").type("Hirthehaven");
    cy.get("#f-county").type("Saskatchewan");
    cy.get("#f-postcode").type("S7R 9F9");
    cy.get("#f-country").type("Canada");
    cy.get("label[for=f-isAirmailRequired]").click();
    cy.get("#f-phoneNumber").type("345");
    cy.get("#f-email").type("test@test.com");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-lpa");
  });
});