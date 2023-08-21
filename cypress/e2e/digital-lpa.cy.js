describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.visit("/lpa/M-1234-9876-4567");
  });

  it("shows case information", () => {
    cy.contains("M-1234-9876-4567");
    cy.get("h1").contains("Zoraida Swanberg");
    cy.get(".govuk-tag.app-tag--draft").contains("Draft");

    cy.contains("1 Complaints");
    cy.contains("2 Investigations");
    cy.contains("3 Tasks");
    cy.contains("4 Warnings");
  });

  it("shows payment information", () => {
    cy.contains("M-1234-9876-4567");
    cy.get("h1").contains("Zoraida Swanberg");

    cy.contains("Fees").click();
    cy.contains("£41.00 expected");
    cy.contains("Record a fee reduction");
  });
});
