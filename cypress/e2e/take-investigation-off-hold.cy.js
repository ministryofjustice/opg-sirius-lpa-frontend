describe("Takes an investigation off hold", () => {
    beforeEach(() => {
        cy.visit("/take-investigation-off-hold?id=175");
    });

    it("takes an investigation off hold", () => {
        cy.contains("Investigation off hold");
        cy.get(".moj-banner").should("not.exist");
        cy.contains("Test title");
        cy.contains("Normal");
        cy.contains("23/01/2022");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
