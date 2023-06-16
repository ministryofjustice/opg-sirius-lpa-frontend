describe("Takes an investigation off hold", () => {
  beforeEach(() => {
    cy.visit("/investigation-hold?id=301");
  });

  it("takes an investigation off hold", () => {
    cy.contains("Take investigation off hold");
    cy.get(".moj-banner").should("not.exist");
    cy.contains("Test title");
    cy.contains("Normal");
    cy.contains("23/01/2022");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
