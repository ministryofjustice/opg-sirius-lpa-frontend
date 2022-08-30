describe("Add a payment", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("OPG-Bypass-Membrane", "1");
        cy.visit("/add-payment?id=800");
    });

    it("adds a payment to the case", () => {
        cy.contains("Add a payment");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#amount").type("41.00");
        cy.get("#source").select("PHONE");
        cy.get("#f-paymentDate").type("2022-03-25");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });

    it("sets the payment date to today", () => {
        cy.clock(Date.UTC(2022, 1, 25), ['Date']); // months in Date starts from 0 so February = 1
        cy.contains("Add a payment");
        cy.contains("7000-0000-0000");
        cy.get(".moj-banner").should("not.exist");
        cy.get('[data-module="select-todays-date"]').click();
        cy.get("#f-paymentDate").should('have.value', "2022-02-25")
    });


});
