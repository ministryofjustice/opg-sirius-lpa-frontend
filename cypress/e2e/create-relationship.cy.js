describe("Create a relationship", () => {
  beforeEach(() => {
    cy.visit("/create-relationship?id=189");
  });

  it("creates a relationship", () => {
    cy.contains("Create Relationship");
    cy.contains("John Doe");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-search").type("7000-0000-0003");
    cy.get(".autocomplete__menu").contains("John Doe (7000-0000-0003)").click();
    cy.get("#f-reason").type("Mother");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
