describe("Delete a payment", () => {
    beforeEach(() => {
        cy.visit("/delete-payment?id=123");
    });

    it("deletes a payment on a case", () => {
        cy.contains("Delete payment");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
