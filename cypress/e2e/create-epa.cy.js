describe("create an EPA", () => {
  it("creates an EPA", () => {
    cy.visit("/create-epa?id=1");

    // step 1
    cy.get("#f-receiptDate").type("2026-06-19");
    cy.get("#f-epaDonorSignatureDate").type("2026-06-19");
    cy.get("#f-epaDonorNoticeGivenDate").type("2026-06-19");
    cy.get("#f-donorHasOtherEpas-2").click();

    // step 2
    cy.get("#f-caseAttorney").click();

    // step 3 tested separately in create-attorney and create-correspondent

    // step 4
    cy.get("#f-paymentByCheque-2").click();
    cy.get("#f-paymentExemption-2").click();

    cy.get("#f-paymentDate").type("2026-06-19");

    cy.contains("button", "Save and exit").click();
  });
});
