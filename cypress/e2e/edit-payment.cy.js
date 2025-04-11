describe("Edit a payment", () => {
  beforeEach(() => {
    cy.visit("/edit-payment?id=123");
  });

  it("edits a payment on the case", () => {
    cy.contains("Edit payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-amount").clear().type("25.50");
    cy.get("#f-paymentDate").type("2022-04-27");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
