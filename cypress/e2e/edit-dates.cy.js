describe("Edit dates", () => {
  beforeEach(() => {
    cy.visit("/edit-dates?id=800&case=lpa");
  });

  it("edits the dates", () => {
    const setSplitDate = (field, day, month, year) => {
      cy.get(`#f-${field}-day`).type(day);
      cy.get(`#f-${field}-month`).type(month);
      cy.get(`#f-${field}-year`).type(year);
    };

    cy.contains("Edit Dates");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");

    setSplitDate("receiptDate", "04", "03", "2022");
    setSplitDate("paymentDate", "04", "03", "2022");
    setSplitDate("filingDate", "04", "03", "2022");
    setSplitDate("dueDate", "04", "03", "2022");
    setSplitDate("registrationDate", "04", "03", "2022");

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
