describe("Change status", () => {
  beforeEach(() => {
    cy.visit("/change-status?id=800&case=lpa");
  });

  it("changes the case status", () => {
    cy.contains("Change status");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-status").select("Perfect");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });

  it("changes the case status and adds a note", () => {
    cy.contains("Change status");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-status").select("Perfect");
    cy.get("#f-notes").type("This is a note");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
