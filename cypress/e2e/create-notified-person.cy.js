describe("Create notified person form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
      },
    });

    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/create-notified-person?id=1&caseId=2");
  });

  it("redirects to lpa form on submit", () => {
    cy.contains("Add a notified person");
    cy.get("#f-salutation").type("Prof");
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").type("Josefina");
    cy.get("#f-surname").type("Vanvolkenburg");
    cy.get(".govuk-details__summary").click();
    cy.get("#f-addressLine1").type("29737 Andrew Plaza");
    cy.get("#f-addressLine2").type("Apt. 814");
    cy.get("#f-addressLine3").type("Gislasonside");
    cy.get("#f-town").type("Hirthehaven");
    cy.get("#f-county").type("Saskatchewan");
    cy.get("#f-postcode").type("S7R 9F9");
    cy.get("#f-country").type("Canada");
    cy.get("#f-noticeGivenDate").type("2023-01-01");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-lpa");
  });

  it("redirects to notified person form when adding another", () => {
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").type("Josefina");
    cy.get("input[type=submit][name=add-another-notified-person]").click();
    cy.url().should("include", "create-notified-person");
    cy.contains("Add a notified person");
  });

  it("does not allow you to add another notified person you are creating the 5th", () => {
    cy.addMock("/lpa-api/v1/cases/4", "GET", {
      status: 200,
      body: {
        id: 4,
        notifiedPersons: [{ id: 11 }, { id: 12 }, { id: 13 }, { id: 14 }],
      },
    });
    cy.visit("/create-notified-person?id=1&caseId=4");
    cy.get("input[type=submit][name=add-another-notified-person]").should(
      "not.exist",
    );
  });
});
