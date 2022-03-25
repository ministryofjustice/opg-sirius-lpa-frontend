describe("Create a task", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("OPG-Bypass-Membrane", "1");
        cy.visit("/create-task?id=800");
    });

    it("creates a task for a user", () => {
        cy.contains("Create Task");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-taskType").select("Check Application");
        cy.get("#f-name").type("A title");
        cy.get("#f-description").type("More words");
        cy.get("#f-assignTo").check();
        cy.get("#f-assigneeUser").type("admin");
        cy.get(".autocomplete__menu").contains("system admin").click();
        cy.get("#f-dueDate").type("9999-03-04");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });

    it("creates a task for a team", () => {
        cy.contains("Create Task");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-taskType").select("Check Application");
        cy.get("#f-name").type("A title");
        cy.get("#f-description").type("More words");
        cy.get("#f-assignTo-2").check();
        cy.get("#f-assigneeTeam").select("Cool Team");
        cy.get("#f-dueDate").type("9999-03-04");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
