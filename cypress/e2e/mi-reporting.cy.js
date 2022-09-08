describe("MI Reporting", () => {
  beforeEach(() => {
    cy.visit("/mi-reporting");
  });

  it("generates a report", () => {
    cy.contains("MI Reporting");
    cy.get("#reportType").select("Number of EPAs received");
    cy.contains("button", "Select").click();

    cy.contains("Number of EPAs received");
    cy.contains("button", "Generate").click();

    cy.contains("Number of EPAs received");
    cy.contains("10");
  });
});
