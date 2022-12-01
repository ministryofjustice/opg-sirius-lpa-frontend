describe("Create a document", () => {
    beforeEach(() => {
        cy.visit("/create-document?id=800&case=lpa");
    });

    it("creates a document on the case", () => {
        cy.contains("Create Draft");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-templateId").type("DD");
        cy.get(".autocomplete__menu").contains("DD: Donor deceased: Blank template").click();
        cy.contains("button", "Select template").click();

        cy.contains("Template: DD");
        cy.get("#f-DD1").click();
        cy.contains("button", "Select inserts").click();
    });
});
