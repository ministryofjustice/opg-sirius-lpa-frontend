describe("Search", () => {
    describe("Searching by name", () => {

        it("finds a person with associated case", () => {
            cy.visit("/search?term=john");
            cy.contains("You searched for: john");
            cy.contains("Showing 1 to 1 of 1 results");
            cy.contains("Donor (1)");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "John Doe");
            $row.should("contain", "123 Somewhere Road");
            $row.should("contain", "perfect");
            $row.should("contain", "LPA - PFA");
            $row
                .contains("7000-8548-8461")
                .should("have.attr", "href")
                .should("contain", "/person/47/23");
        });
    });

    describe("Search features", () => {
        beforeEach(() => {
            cy.visit("/search?term=abcdefg");
        });

        it("it shows/hides filter panel", () => {
            cy.contains(".govuk-button", "Hide filters").click();
            cy.contains("Apply filters").should("not.be.visible");
            cy.contains(".govuk-button", "Show filters").click();
            cy.contains("Apply filters").should("be.visible");
        });

        it("enables the person type filters on selection", () => {
            cy.contains(".govuk-checkboxes__item", "Attorney").find("input").check();
            cy.contains(".govuk-checkboxes__item", "Trust corporation").find("input").check();
            cy.get("button[type=submit]").click();
            cy.contains(".moj-filter__tag", "Attorney");
            cy.contains(".moj-filter__tag", "Trust Corporation");
        });

        it("can clear all filters", () => {
            cy.contains(".govuk-checkboxes__item", "Donor").find("input").check();
            cy.contains(".govuk-checkboxes__item", "Attorney").find("input").check();
            cy.contains(".govuk-checkboxes__item", "Client").find("input").check();
            cy.get("button[type=submit]").click();
            cy.contains(".moj-filter__tag", "Donor");
            cy.contains(".moj-filter__tag", "Attorney");
            cy.contains(".moj-filter__tag", "Client");
            cy.contains(".moj-filter__selected-heading", "Clear filters").find("a").click();
            cy.get('.moj-filter__tag').should('not.exist');
        });
    });
});
