describe("Edit investigation", () => {
  beforeEach(() => {
    cy.visit("/edit-investigation?id=300");
  });

  it("edits a investigation", () => {
    cy.contains("Edit Investigation");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-title").should("have.value", "Test title");
    cy.get("#f-information").should("have.value", "Some test info");
    cy.contains(".govuk-radios__label", "Normal")
        .parent()
        .get("input")
        .should("be.checked");
    cy.get("#f-dateReceived").should("have.value", "2022-01-23");
    cy.get("#f-approvalDate").type("2022-04-05");
    cy.get("#f-riskAssessmentDate").type("2022-04-05");
    cy.get("#f-approvalOutcome").select("Court Application");
    cy.get("#f-investigationClosureDate").type("2022-04-05");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
