describe("Assign task", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("OPG-Bypass-Membrane", "1");
        cy.visit("/assign-task?id=990");
    });

    it("assigns a task", () => {
        cy.contains("Assign Task");
        cy.get(".moj-banner").should("not.exist");
        cy.contains(".govuk-radios__item", "User").find("input").check();
        cy.get("#f-assigneeUser").type("admin");
        cy.get(".autocomplete__menu").contains("system admin").click();
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
