describe("Update correspondent", () => {
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

  it("has a back link to the EPA form", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: { id: 2 },
    });

    cy.visit("/edit-correspondent?id=1&caseId=2&caseType=epa");
    cy.get(".govuk-back-link")
      .should("exist")
      .and("have.attr", "href")
      .and(
        "include",
        "/create-epa?id=1&caseId=2#accordion-create-epa-heading-3",
      );
  });

  it("updates a correspondent on an EPA", () => {
    cy.addMock("/lpa-api/v1/epas/2", "PUT", {
      status: 200,
      body: {},
    });

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
    cy.contains("Edit EPA");
    cy.contains("Melanie Vanvolkenburg");
    cy.get("#f-update-correspondent-3 .govuk-visually-hidden").should(
      "contain.text",
      "correspondent Melanie Vanvolkenburg",
    );
    cy.get("#f-update-correspondent-3").click();

    cy.contains("Update correspondent details");
    cy.get("#f-firstname").should("have.value", "Melanie");
    cy.get("#f-surname").should("have.value", "Vanvolkenburg");
    cy.get("input[type=submit][name=add-another]").should("not.exist");

    cy.get("#f-firstname").clear().type("Mindy");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-epa");
  });

  it("updates a correspondent on an LPA", () => {
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

    cy.visit("/edit-correspondent?id=1&caseId=2&caseType=lpa");

    cy.contains("Update correspondent details");
    cy.get("#f-firstname").should("have.value", "Melanie");
    cy.get("#f-surname").should("have.value", "Vanvolkenburg");
    cy.get("input[type=submit][name=add-another]").should("not.exist");

    cy.get("#f-firstname").clear().type("Mindy");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-lpa");
  });
});
