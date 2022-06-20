describe("Add a complaint", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("OPG-Bypass-Membrane", "1");
        cy.visit("/add-complaint?id=800&case=lpa");
    });

    it("adds a complaint to the case", () => {
        cy.contains("Add Complaint");
        cy.contains("LPA 7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.contains(".govuk-radios__label", "Major").click();
        cy.get("#f-summary").type("A title");
        cy.get("#f-description").type("A description");
        cy.get("#f-receivedDate").type("2022-04-05");
        cy.contains(".govuk-radios__label", "OPG Decisions").click();
        cy.get("#f-subCategory-02").select("Fee Decision");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
