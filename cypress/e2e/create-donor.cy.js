describe("Create donor", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/create-donor");
  });

  it("creates a donor", () => {
    cy.contains("Create Donor");
    cy.get(".moj-banner").should("not.exist");

    cy.get("#f-salutation").type("Prof");
    cy.get("#f-firstname").type("Melanie");
    cy.get("#f-middlenames").type("Josefina");
    cy.get("#f-surname").type("Vanvolkenburg");
    cy.get("#f-dob").type("1978-04-19");
    cy.get("#f-previousNames").type("Colton Bacman");
    cy.get("#f-otherNames").type("Mel");
    cy.get("#f-addressLine1").type("29737 Andrew Plaza");
    cy.get("#f-addressLine2").type("Apt. 814");
    cy.get("#f-addressLine3").type("Gislasonside");
    cy.get("#f-town").type("Hirthehaven");
    cy.get("#f-county").type("Saskatchewan");
    cy.get("#f-postcode").type("S7R 9F9");
    cy.get("#f-country").type("Canada");
    cy.get("#f-isAirmailRequired").click();
    cy.get("#f-phoneNumber").type("072345678");
    cy.get("#f-email").type("m.vancolkenburg@ca.test");
    cy.get("#f-correspondenceBy-email").click();
    cy.get("#f-correspondenceBy-phone").click();
    cy.get("#f-researchOptOut").click();

    cy.get("button[type=submit]").click();
    cy.get(".govuk-notification-banner").should("exist");
    cy.get(".govuk-notification-banner").contains(
      "Person 7000-0290-0192 was created"
    );
  });
});
