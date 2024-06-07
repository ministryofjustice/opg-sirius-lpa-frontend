describe("Clear task on a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/tasks/990/mark-as-completed", "PUT", {
      status: 200,
    });

    cy.visit("/clear-task?id=990");
  });

  it("marks a task as completed", () => {
    cy.contains("Clear Task");
    cy.get(".moj-banner").should("not.exist");
    cy.contains("Task:");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
