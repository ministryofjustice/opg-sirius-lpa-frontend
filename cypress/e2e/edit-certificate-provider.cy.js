describe("Edit certificate provider form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/33", "GET", {
      status: 200,
      body: {
        id: 33,
        caseId: 2,
        salutation: "Dr",
        firstName: "Test",
        middleNames: "Somebody",
        surname: "Testing",
        addressLine1: "Flat 1",
        addressLine2: "3 Church Road",
        addressLine3: "",
        town: "Blackpool",
        county: "Lancashire",
        postcode: "FY48 7CY",
        country: "United Kingdom",
      },
    });

    cy.visit("/edit-certificate-provider?id=1&caseId=2&personId=33");
  });

  it("populates certificate provider details ", () => {
    cy.contains("Edit a certificate provider");
    cy.get("#f-salutation").should("have.value", "Dr");
    cy.get("#f-firstname").should("have.value", "Test");
    cy.get("#f-middlenames").should("have.value", "Somebody");
    cy.get("#f-surname").should("have.value", "Testing");
    cy.get("#f-addressLine1").should("have.value", "Flat 1");
    cy.get("#f-addressLine2").should("have.value", "3 Church Road");
    cy.get("#f-addressLine3").should("have.value", "");
    cy.get("#f-town").should("have.value", "Blackpool");
    cy.get("#f-county").should("have.value", "Lancashire");
    cy.get("#f-postcode").should("have.value", "FY48 7CY");
    cy.get("#f-country").should("have.value", "United Kingdom");
  });

  it("redirects to lpa form on submit", () => {
    cy.addMock("/lpa-api/v1/certificate-providers/33", "PUT", {
      status: 200,
      body: {},
    });

    cy.contains("Edit a certificate provider");
    cy.get("#f-salutation").clear();
    cy.get("#f-salutation").type("Prof");

    cy.get("#f-firstname").clear();
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").clear();
    cy.get("#f-middlenames").type("Josefina");
    cy.get("#f-surname").clear();
    cy.get("#f-surname").type("Vanvolkenburg");

    cy.get("#f-addressLine1").clear();
    cy.get("#f-addressLine1").type("29737 Andrew Plaza");
    cy.get("#f-addressLine2").clear();
    cy.get("#f-addressLine2").type("Apt. 814");
    cy.get("#f-addressLine3").type("Gislasonside");

    cy.get("#f-town").clear();
    cy.get("#f-town").type("Hirthehaven");
    cy.get("#f-county").clear();
    cy.get("#f-county").type("Saskatchewan");
    cy.get("#f-postcode").clear();
    cy.get("#f-postcode").type("S7R 9F9");
    cy.get("#f-country").clear();
    cy.get("#f-country").type("Canada");

    cy.get("button[type=submit]").click();
    cy.url().should("include", "create-lpa");
  });

  it("has a back link to the LPA form", () => {
    cy.get(".govuk-back-link")
      .should("exist")
      .and("have.attr", "href")
      .and(
        "include",
        "/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-3",
      );
  });
});
