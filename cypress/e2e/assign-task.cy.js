describe("Assign task", () => {
  beforeEach(() => {
    cy.visit("/assign-task?id=990");
  });

  it("assigns a task", () => {
    cy.contains("Assign Task");
    cy.get(".moj-banner").should("not.exist");
    cy.contains("label", "User").click();
    cy.get("#f-assigneeUser").type("admin");
    cy.get(".autocomplete__menu").contains("system admin").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
