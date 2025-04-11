describe("Delete a relationship", () => {
  beforeEach(() => {
    cy.visit("/delete-relationship?id=189");
  });

  it("deletes a relationship", () => {
    cy.contains("Delete Relationship");
    cy.contains("John Doe");
    cy.get(".moj-alert").should("not.exist");
    cy.contains("label", "John Doe").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
