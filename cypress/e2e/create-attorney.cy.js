const fillInAttorneyDetails = () => {
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
};

describe("Create or Update Attorney", () => {
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
  });

  it("creates an attorney on an EPA", () => {
    fillInAttorneyDetails();
    cy.contains("Add an attorney");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-epa");
  });

  it("creates an attorney on an EPA and add another attorney", () => {
    fillInAttorneyDetails();
    cy.contains("Add an attorney");
    cy.get("input[type=submit][name=add-another]").click();
    cy.url().should("include", "create-attorney");
  });

  it("updates an existing attorney on an EPA", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        attorneys: [
          {
            id: 3,
            firstname: "Rudolph",
            surname: "Stotesbury",
            relationshipToDonor: "NO RELATION",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/attorneys/3", "PUT", {
      status: 200,
      body: {},
    });

    cy.visit("/create-epa?id=2&caseId=2");
    cy.contains("Edit EPA");
    cy.contains("Rudolph Stotesbury");
    cy.get("#f-update-attorney-3 .govuk-visually-hidden").should(
      "contain.text",
      "attorney Rudolph Stotesbury",
    );
    cy.get("#f-update-attorney-3")
      .should(
        "have.attr",
        "href",
        "/create-attorney?id=2&caseId=2&attorneyId=3",
      )
      .click();

    cy.contains("Update attorney details");
    cy.get("#f-firstname").should("have.value", "Rudolph");
    cy.get("#f-surname").should("have.value", "Stotesbury");
    cy.get("#f-relationshipToDonor").should("have.value", "NO RELATION");
    cy.get("input[type=submit][name=add-another]").should("not.exist");

    cy.get("#f-firstname").clear().type("Rafael");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-epa");
  });
});
