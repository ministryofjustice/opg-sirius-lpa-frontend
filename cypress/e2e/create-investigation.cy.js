describe("Add an investigation", () => {
    beforeEach(() => {
        cy.visit("/create-investigation?id=800&case=lpa");
    });

    it("creates an investigation on the case", () => {
        cy.contains("Create Investigation");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-title").type("A title");
        cy.get("#f-information").type("Some information");
        cy.contains(".govuk-radios__label", "Priority").click();
        cy.get("#f-dateReceived").type("2022-03-04");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
