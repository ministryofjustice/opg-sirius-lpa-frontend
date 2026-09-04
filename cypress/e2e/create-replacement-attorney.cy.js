const fillInAttorneyDetails = () => {
  cy.get("#f-salutation").type("Prof");
  cy.get("#f-firstname").type("Melanie");
  cy.get("#f-middlenames").type("Josefina");
  cy.get("#f-surname").type("Vanvolkenburg");
  cy.get("#f-dob").type("1978-04-19");
  cy.get(".govuk-details__summary").click();
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
};

describe("Create or Update Replacement Attorney on an LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: [],
    });

    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        caseSubtype: "pfa",
        receiptDate: "19/06/2026",
        replacementAttorneys: [
          {
            id: 3,
            firstname: "Rudolph",
            surname: "Stotesbury",
          },
        ],
      },
    });

    cy.visit("/create-replacement-attorney?id=1&caseId=2");
  });

  it("creates a replacement attorney on an LPA", () => {
    fillInAttorneyDetails();
    cy.contains("Add a replacement attorney");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-lpa");
  });

  it("creates a replacement attorney on an LPA and add another replacement attorney", () => {
    fillInAttorneyDetails();
    cy.contains("Add a replacement attorney");
    cy.get("input[type=submit][name=add-another]").click();
    cy.url().should("include", "create-replacement-attorney");
  });

  it("has a back link to the LPA form", () => {
    cy.get(".govuk-back-link")
      .should("exist")
      .and("have.attr", "href")
      .and(
        "include",
        "/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-1",
      );
  });

  it("updates an existing replacement attorney on an LPA", () => {
    cy.addMock("/lpa-api/v1/replacement-attorneys/3", "PUT", {
      status: 200,
      body: {},
    });

    cy.visit("/create-replacement-attorney?id=1&caseId=2&attorneyId=3");
    cy.contains("Update replacement attorney details");
    cy.get("#f-firstname").should("have.value", "Rudolph");
    cy.get("#f-surname").should("have.value", "Stotesbury");
    cy.get("input[type=submit][name=add-another]").should("not.exist");

    cy.get("#f-firstname").clear().type("Rafael");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-lpa");
  });

  it("should show the trust corporation link on a pfa lpa", () => {
    cy.contains("Add a trust corporation as a replacement attorney");
  });

  it("should not show the trust corporation link on a hw lpa", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        caseSubtype: "hw",
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

    cy.visit("/create-attorney?id=1&caseId=2&caseType=lpa");

    cy.contains("Add a trust corporation as a replacement attorney").should(
      "not.exist",
    );
  });
});
