describe("Add a fee decline decision to a non-digital LPA", () => {
  beforeEach(() => {
    cy.visit("/create-decline-fee-decision?id=801");
    cy.addMock("/lpa-api/v1/cases/801/fee-decisions", "GET", {
      status: 200,
      body: [],
    });
  });

  it("adds a fee decline decision to the case", () => {
    cy.contains("Decline a fee reduction");
    cy.contains("7000-0000-0001");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-decisionType").select("Declined exemption");
    cy.get("#f-decisionReason").type("Invalid evidence");
    cy.get("#f-decisionDate").type("2023-10-09");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Fee decision added");
  });
});
