describe("MI Reporting", () => {
  beforeEach(() => {
    cy.visit("/mi-reporting");
  });

  it("generates a report", () => {
    cy.contains("MI Reporting");
    cy.get("#f-reportType").select("Number of EPAs received");
    cy.contains("button", "Select").click();

    cy.contains("Number of EPAs received");
    cy.contains("button", "Generate").click();

    cy.contains("Number of EPAs received");
    cy.contains("10");

    // Make the link open in a new tab because otherwise Cypress
    // throws cross-domain errors
    cy.contains("a", "Download").then(($input) => {
      $input[0].setAttribute("target", "_blank");
    });

    cy.contains("a", "Download").click();
    cy.contains("a", "Download")
      .invoke("attr", "class")
      .should("contain", "govuk-button--disabled");
    cy.contains("Your download will open in a new window when ready");
  });
});
