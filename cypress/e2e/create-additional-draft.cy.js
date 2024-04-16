describe("Create Additional Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });

    cy.visit("/create-additional-draft-lpa?id=990");
  });

  it("creates an additional digital LPA", () => {
    cy.contains("Create a draft LPA for Steven Munnell");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("The donor, using the address above").click();

    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-panel").contains(
      "2 draft LPAs for Steven Munnell have been created.",
    );
  });
});
