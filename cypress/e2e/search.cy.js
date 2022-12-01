describe("Search", () => {
    beforeEach(() => {
        cy.visit("/search?term=bob");
    });

    it("finds a person when searched", () => {
        cy.contains("Search results");
        cy.contains("You searched for: bob");
        cy.contains("Showing 1 to 2 of 2 cases");
        cy.contains("Bob Smith");
    });
});

describe("Search", () => {
    beforeEach(() => {
        cy.visit("/search?term=abcdefg");
    });

    it("it cannot find any results", () => {
        cy.contains("Search results");
        cy.contains("No cases were found");
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