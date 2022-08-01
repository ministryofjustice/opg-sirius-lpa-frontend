describe("Mi Reporting", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/mi-reporting");
  });

  it("generates a report", () => {
    cy.contains("Mi Reporting");
    cy.get("#reportType").select("Number of EPAs received");
    cy.contains("button", "Select").click();

    cy.contains("Mi Reporting: Number of EPAs received");
    cy.contains("button", "Generate").click();

    cy.contains("Mi Reporting: Number of EPAs received");
    cy.contains("10");
  });
});
