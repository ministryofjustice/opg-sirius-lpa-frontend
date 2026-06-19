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

    cy.get("button[type=submit]").click();
  });

  it("updates an existing EPA", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        receiptDate: "19/06/2026",
        epaDonorSignatureDate: "19/06/2026",
        epaDonorNoticeGivenDate: "19/06/2026",
        donorHasOtherEpas: false,
        caseAttorneyJointly: true,
        paymentByCheque: false,
        paymentExemption: false,
        paymentDate: "19/06/2026",
      },
    });

    cy.visit("/create-epa?id=1&caseId=2");

    // step 1
    cy.get("#f-receiptDate").should("have.value", "2026-06-19");
    cy.get("#f-epaDonorSignatureDate").should("have.value", "2026-06-19");
    cy.get("#f-epaDonorNoticeGivenDate").should("have.value", "2026-06-19");
    cy.get("#f-donorHasOtherEpas-2").should("be.checked");

    // step 2
    cy.get("#f-caseAttorney-3").should("be.checked");

    // step 3 tested separately in create-attorney and create-correspondent

    // step 4
    cy.get("#f-paymentByCheque-2").should("be.checked");
    cy.get("#f-paymentExemption-2").should("be.checked");

    cy.get("#f-paymentDate").should("have.value", "2026-06-19");
    cy.get("#f-paymentDate").type("2026-06-18");

    cy.get("button[type=submit]").click();
  });
});
