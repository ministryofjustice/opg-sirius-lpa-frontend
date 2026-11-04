describe("Delete a fee reduction", () => {
  beforeEach(() => {
    cy.visit("/delete-fee-reduction?id=124");
    cy.addMock("/lpa-api/v1/cases/802/fee-decisions", "GET", {
      status: 200,
      body: [],
    });
  });

  it("deletes a fee reduction on a case", () => {
    cy.contains("Delete remission");
    cy.contains("7000-0000-0002");
    cy.get(".moj-alert").should("not.exist");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
