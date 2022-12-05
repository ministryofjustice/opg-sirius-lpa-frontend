describe("Search", () => {
    describe("Searching by name", () => {
        beforeEach(() => {
            cy.visit("/search?term=bob");
        });

        it("finds a person when searching by name", () => {
            cy.contains("Search results");
            cy.contains("You searched for: bob");
            cy.contains("Showing 1 to 2 of 2 cases");
            cy.contains("Bob Smith");
        });

        it("it cannot find any results", () => {
            cy.visit("/search?term=someone");
            cy.contains("Search results");
            cy.contains("No cases were found");
        });
    });

    describe("Search deleted case", () => {
        beforeEach(() => {
            cy.visit("/search?term=700000005555");
        });

        it("finds a deleted case when searching by uid", () => {
            cy.contains("Search results");
            cy.contains("7000-0000-5555");
            cy.contains("A54123456789");
            cy.contains("LPA");
            cy.contains("02/12/2022");
            cy.contains("Return - unpaid");
            cy.contains("LPA was not paid for after 12 months");
        });
    });


    describe("Search features", () => {
        beforeEach(() => {
            cy.visit("/search?term=abcdefg");
        });

        it("it shows/hides filter panel", () => {
            cy.contains(".govuk-button", "Hide filter").click();
            cy.contains("Apply filters").should("not.be.visible");
            cy.contains(".govuk-button", "Show filter").click();
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