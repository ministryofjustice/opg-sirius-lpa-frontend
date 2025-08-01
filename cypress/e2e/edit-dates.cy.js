describe("Edit dates", () => {
  beforeEach(() => {
    cy.visit("/edit-dates?id=800&case=lpa");
  });

  it("edits the dates", () => {
    cy.contains("Edit Dates");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-receiptDate").type("2022-03-04");
    cy.get("#f-paymentDate").type("2022-03-04");
    cy.get("#f-dueDate").type("2022-03-04");
    cy.get("#f-registrationDate").type("2022-03-04");
    cy.get("#f-dispatchDate").type("2022-03-04");
    cy.get("#f-cancellationDate").type("2022-03-04");
    cy.get("#f-rejectedDate").type("2022-03-04");
    cy.get("#f-invalidDate").type("2022-03-04");
    cy.get("#f-withdrawnDate").type("2022-03-04");
    cy.get("#f-revokedDate").type("2022-03-04");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
