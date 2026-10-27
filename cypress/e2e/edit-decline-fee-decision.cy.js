describe("Edit a declined fee reduction", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/801", "GET", {
      status: 200,
      body: {
        id: 801,
        uid: "7000-0000-0001",
        caseType: "LPA",
        status: "Registered",
        expectedPaymentTotal: 8200,
      },
    });

    cy.addMock("/lpa-api/v1/cases/801/fee-decisions", "GET", {
      status: 200,
      body: [
        {
          id: 123,
          decisionDate: "09/10/2023",
          decisionType: "DECLINED_EXEMPTION",
          decisionReason: "Invalid evidence",
        },
      ],
    });

    cy.addMock("/lpa-api/v1/cases/801/payments", "GET", {
      status: 200,
      body: [],
    });

    cy.visit("/payments/801");
  });

  it("opens the existing decision, prepopulates the form, and saves edits", () => {
    cy.get("#f-decline-fee-reductions-tab").click();
    cy.contains("a", "Edit declined fee reduction").click();

    cy.contains("h1", "Edit a declined fee reduction");
    cy.get("#f-decisionType").should("have.value", "DECLINED_EXEMPTION");
    cy.get("#f-decisionReason").should("have.value", "Invalid evidence");
    cy.get("#f-decisionDate").should("have.value", "2023-10-09");

    cy.addMock("/lpa-api/v1/cases/801/fee-decisions/123", "PUT", {
      status: 200,
      body: {},
    });

    cy.get("#f-decisionReason").clear().type("Updated evidence");
    cy.get("#f-decisionDate").clear().type("2023-10-10");
    cy.get("button[type=submit]").click();

    cy.url().should("contain", "/payments/801");
    cy.get(".moj-alert").should("contain", "Decline a fee reduction updated");
  });
});
