describe("View a payment", () => {
    describe("No payments on case", () => {
        it("displays default message when there are no payments on the case", () => {
            cy.visit("/payments?id=999");
            cy.contains("7000-0000-0001");
            cy.contains("There is currently no fee data available to display.");
        });
    });

    describe("Payments on case", () => {
        beforeEach(() => {
            cy.visit("/payments?id=800");
        });

        it("displays payment information if there is a payment on the case", () => {
            cy.contains("7000-0000-0000");
            cy.contains("Total paid");
            cy.contains("£41.00");
            cy.contains("Fee details");
            cy.contains("Payments");
            cy.get('.govuk-details__summary-text').click();
            cy.contains("Amount");
            cy.contains("£41.00");
            cy.contains("Date of payment:");
            cy.contains("2022-01-23");
            cy.contains("Method");
            cy.get(".govuk-link").contains("Edit payment");
        });

        it("displays add payment and apply fee reduction buttons", () => {
            cy.get(".govuk-button").contains("Add payment");
            cy.get(".govuk-button").contains("Apply fee reduction");
        });
    });
});
