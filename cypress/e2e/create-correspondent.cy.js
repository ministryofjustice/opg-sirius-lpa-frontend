describe("Create or update correspondent", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });

    cy.addMock("/lpa-api/v1/persons/3", "PUT", {
      status: 200,
      body: {},
    });
  });

  it("creates a correspondent on an EPA", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
      },
    });

    cy.visit("/create-correspondent?id=1&caseId=2");

    cy.contains("Add a correspondent");
    cy.get("#f-salutation").type("Prof");
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").type("Josefina");
    cy.get("#f-surname").type("Vanvolkenburg");
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

  it("updates a correspondent on an EPA", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        correspondent: {
            id: 3,
            firstName: "Melanie",
            surname: "Vanvolkenburg",
        },
      },
    });

    cy.visit("/create-epa?id=2&caseId=2");
    cy.contains("Edit EPA")
    cy.contains("Melanie Vanvolkenburg");
      cy.get("#f-update-correspondent-3 .govuk-visually-hidden").should(
      "contain.text",
      "correspondent Melanie Vanvolkenburg",
    );
    cy.get("#f-update-correspondent-3").should(
      "have.attr",
      "href",
      "/create-correspondent?id=2&caseId=2",
    ).click();

    cy.contains("Update correspondent details");
    cy.get("#f-firstname").should("have.value", "Melanie");
    cy.get("#f-surname").should("have.value", "Vanvolkenburg");
    cy.get("input[type=submit][name=add-another]").should("not.exist");

    cy.get('#f-firstname').clear().type("Mindy");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-epa");
  });
});
