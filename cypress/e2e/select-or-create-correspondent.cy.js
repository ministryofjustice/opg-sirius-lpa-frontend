describe("Select or create correspondent", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        attorneys: [
          {
            id: 3,
            firstname: "Melanie",
            surname: "Vanvolkenburg",
          },
          {
            id: 4,
            firstname: "Will",
            surname: "Niesborella",
          },
        ],
        donor: {
          id: 5,
          firstname: "Sarah",
          surname: "Jones",
        },
        replacementAttorneys: [
          {
            id: 6,
            firstname: "Steven",
            surname: "Munnell",
          },
        ],
        certificateProviders: [
          {
            id: 7,
            firstname: "Will",
            surname: "Oswald",
          },
        ],
        notifiedPersons: [
          {
            id: 8,
            firstname: "James",
            surname: "Rubin",
          },
        ],
        trustCorporations: [
          {
            id: 9,
            companyName: "ACME",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });
  });

  it("can select an existing attorney to create a correspondent", () => {
    cy.visit("/select-or-create-correspondent?id=1&caseId=2&caseType=epa");
    cy.contains("Add a correspondent");
    cy.get("label[for=f-actor-1]").should(
      "contain.text",
      "Melanie Vanvolkenburg",
    );
    cy.get("label[for=f-actor-2]")
      .should("contain.text", "Will Niesborella")
      .click();
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-epa");
  });

  it("displays the correct actors for an lpa", () => {
    cy.visit("/select-or-create-correspondent?id=1&caseId=2&caseType=lpa");
    cy.contains("Add a correspondent");
    cy.get("label[for=f-actor-donor]").should("contain.text", "Sarah Jones");
    cy.get("label[for=f-attorney-1]").should(
      "contain.text",
      "Melanie Vanvolkenburg",
    );
    cy.get("label[for=f-replacement-1]").should(
      "contain.text",
      "Steven Munnell",
    );
    cy.get("label[for=f-certified-provider-1]").should(
      "contain.text",
      "Will Oswald",
    );
    cy.get("label[for=f-notified-person-1]").should(
      "contain.text",
      "James Rubin",
    );
    cy.get("label[for=f-trust-corporation-1]").should("contain.text", "ACME");
  });

  it("select create a new correspondent by default", () => {
    cy.visit("/select-or-create-correspondent?id=1&caseId=2&caseType=epa");
    cy.contains("Add a correspondent");
    cy.get("input#f-actor-new").should("be.checked");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-correspondent");
  });
});
