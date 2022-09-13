describe("Edit a payment", () => {
    beforeEach(() => {
        cy.visit("/edit-payment?payment=123");
    });

    it("edits a payment on the case", () => {
        cy.contains("Edit payment");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-amount").clear().type("45.50");
        cy.get("#f-source").select("ONLINE");
        cy.get("#f-paymentDate").type("2022-07-16");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
