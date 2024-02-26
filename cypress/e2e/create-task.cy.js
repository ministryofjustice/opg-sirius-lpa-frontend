describe("Create a task", () => {
  beforeEach(() => {
    cy.visit("/create-task?id=800");
  });

  it("creates a task for a user", () => {
    cy.contains("Create Task");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("Something");
    cy.get("#f-description").type("More words");
    cy.contains("label", "User").click();
    cy.get("#f-assigneeUser").type("admin");
    cy.get(".autocomplete__menu").contains("system admin").click();
    cy.get("#f-dueDate").type("2035-03-04");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });

  it("creates a task for a team", () => {
    cy.contains("Create Task");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("A title");
    cy.get("#f-description").type("A description");
    cy.contains("label", "Team").click();
    cy.get("#f-assigneeTeam").select("Cool Team");
    cy.get("#f-dueDate").type("2035-03-04");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
