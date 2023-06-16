describe("Delete a fee reduction", () => {
  beforeEach(() => {
    cy.visit("/delete-fee-reduction?id=124");
  });

  it("deletes a fee reduction on a case", () => {
    cy.contains("Delete remission");
    cy.contains("7000-0000-0002");
    cy.get(".moj-banner").should("not.exist");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
