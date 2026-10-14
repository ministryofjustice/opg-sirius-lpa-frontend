describe("Delete a payment", () => {
  beforeEach(() => {
    cy.visit("/delete-payment?id=123");
    cy.addMock("/lpa-api/v1/cases/800/fee-decisions", "GET", {
      status: 200,
      body: [],
    });
  });

  it("deletes a payment on a case", () => {
    cy.contains("Delete payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
