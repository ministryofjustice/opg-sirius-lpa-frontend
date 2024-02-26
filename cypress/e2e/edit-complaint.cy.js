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
    cy.get("#f-complainantName").type("Someones name");
    cy.get("#f-summary").should("have.value", "This and that");
    cy.get("#f-description").should("have.value", "This is seriously bad");
    cy.get("#f-receivedDate").should("have.value", "2022-04-05");
    cy.contains("label", "OPG Decisions")
      .parent()
      .get("input")
      .should("be.checked");
    cy.contains(".govuk-radios__label", "OPG Decisions").click();
    cy.get("label[for=f-compensation-type-0]").click();
    cy.get("#f-compensation-amount-0").type("150.00");
    cy.get("#f-subCategory-02").should("have.value", "18");
    cy.get("#f-complainantCategory").select("LPA Donor");
    cy.get("#f-origin").select("Phone call");
    cy.get("#f-resolutionDate").type("2022-06-07");
    cy.contains(".govuk-radios__label", "Complaint Upheld").click();
    cy.get("#f-resolutionInfo").type("Because...");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
