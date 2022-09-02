describe("Edit a payment", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("OPG-Bypass-Membrane", "1");
        cy.visit("/edit-payment?id=800&payment=123");
    });

    it("edits a payment on the case", () => {
        cy.contains("Edit a payment");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-amount").type("45.50");
        cy.get("#f-source").select("ONLINE");
        cy.get("#f-paymentDate").type("2022-07-16");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
