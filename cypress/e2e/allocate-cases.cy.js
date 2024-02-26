describe("Allocate cases", () => {
  beforeEach(() => {
    cy.visit("/allocate-cases?id=800");
  });

  it("allocates a case", () => {
    cy.contains("Allocate Case");
    cy.get(".moj-banner").should("not.exist");
    cy.contains("label", "User").click();
    cy.get("#f-assigneeUser").type("admin");
    cy.get(".autocomplete__menu").contains("system admin").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
