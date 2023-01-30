describe("Edit complaint", () => {
  beforeEach(() => {
    cy.visit("/edit-complaint?id=986");
  });

  it("edits a complaint", () => {
    cy.contains("Edit Complaint");
    cy.get(".moj-banner").should("not.exist");
    cy.contains(".govuk-radios__label", "Major")
      .parent()
      .get("input")
      .should("be.checked");
    cy.get("#f-investigatingOfficer").should("have.value", "Test Officer");
    cy.get("#f-summary").should("have.value", "This and that");
    cy.get("#f-description").should("have.value", "This is seriously bad");
    cy.get("#f-receivedDate").should("have.value", "2022-04-05");
    cy.contains(".govuk-radios__label", "Correspondence")
      .parent()
      .get("input")
      .should("be.checked");
    cy.get("#f-subCategory-01").should("have.value", "07");
    cy.get("#f-complainantCategory").select("LPA Donor");
    cy.get("#f-origin").select("Phone call");
    cy.get("#f-resolutionDate").type("2022-05-06");
    cy.contains(".govuk-radios__label", "Complaint Not Upheld").click();
    cy.get("#f-resolutionInfo").type("Because...");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
