describe("Delete a declined fee decision form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/fee-decisions/123", "GET", {
      status: 200,
      body: {
        id: 123,
        decisionDate: "03/02/2026",
        decisionReason: "Test",
        decisionType: "DECLINED_EXEMPTION",
      },
    });

    cy.addMock("/lpa-api/v1/fee-decisions/123", "DELETE", {
      status: 204,
      body: [],
    });

    cy.addMock("/lpa-api/v1/cases/801/fee-decisions", "GET", {
      status: 200,
      body: [],
    });
    cy.visit("/delete-decline-fee-decision?caseId=801&feeDecisionId=123");
  });

  it("deletes a payment on a case", () => {
    cy.contains("Delete declined exemption");
    cy.contains("Are you sure you want to delete the declined exemption?");
    cy.get(".moj-alert").should("not.exist");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
