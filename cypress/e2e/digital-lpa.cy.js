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
  });

  it("shows document information", () => {
    cy.contains("M-1234-9876-4567");
    cy.get("h1").contains("Zoraida Swanberg");
    cy.contains("Documents").click();

    cy.contains("Mr Test Person - Blank Template");
    cy.contains("[OUT]");
    cy.contains("24 August 2023");
    cy.contains("EP-BB");

    cy.contains("John Doe - Donor deceased: Case Withdrawn");
    cy.contains("[OUT]");
    cy.contains("15 May 2023");
    cy.contains("DD-4");
  });

  it("shows task list", () => {
    cy.contains("M-1234-9876-4567");
    cy.get("h2").contains("Tasks");
    cy.get("ul[data-role=tasks-list] li").should((elts) => {
      expect(elts).to.have.length(3);
      expect(elts).to.contain("Review reduced fee eligibility (Super Team)");
      expect(elts).to.contain(
        "Review application correspondence (Marvellous Team)",
      );
      expect(elts).to.contain("Another task (Super Team)");
    });
  });

  it("redirects to create a task via case actions", () => {
    cy.get("select#case-actions").select('Create a task');
    cy.location('pathname').should('eq', '/create-task');
  });

  it("redirects to create a warning via case actions", () => {
    cy.get("select#case-actions").select('Create a warning');
    cy.location('pathname').should('eq', '/create-warning');
  });
});
