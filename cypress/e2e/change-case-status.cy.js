describe("Change case status", () => {
  beforeEach(() => {
    cy.visit("/change-case-status?uid=M-1234-9876-4567");
  });

  it("changes the digital lpa case status", () => {
    cy.contains("Change case status");
    cy.contains("M-1234-9876-4567");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-status").select("In progress");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
