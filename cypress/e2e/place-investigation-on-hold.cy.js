describe("Puts an investigation on hold", () => {
    beforeEach(() => {
        cy.visit("/place-investigation-on-hold?id=300");
    });

    it("places an investigation on hold", () => {
        cy.contains("Investigation on hold");
        cy.get(".moj-banner").should("not.exist");
        cy.contains("Test title");
        cy.contains("Normal");
        cy.contains("23/01/2022");
        cy.contains(".govuk-radios__label", "Police Investigation").click();
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
