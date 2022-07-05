describe("Edit dates", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/edit-dates?id=800&case=lpa");
  });

  it("edits the dates", () => {
    cy.contains("Edit Dates");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-receiptDate").type("2022-03-04");
    cy.get("#f-dueDate").type("2022-03-04");
    cy.get("#f-registrationDate").type("2022-03-04");
    cy.get("#f-dispatchDate").type("2022-03-04");
    cy.get("#f-cancellationDate").type("2022-03-04");
    cy.get("#f-rejectedDate").type("2022-03-04");
    cy.get("#f-invalidDate").type("2022-03-04");
    cy.get("#f-withdrawnDate").type("2022-03-04");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
