describe("Search", () => {
    describe("Searching by name", () => {
        beforeEach(() => {
            cy.visit("/search", {
                qs: {'term': 'bob'}
            });
        });

        it("finds a person when searching by name", () => {
            cy.contains("You searched for: bob");
            cy.contains("Showing 1 to 1 of 1 cases");
            cy.contains("Donor (1)");
            const $row = cy.get("table > tbody > tr");
            $row
                .contains("7000-8548-8461")
                .should("have.attr", "href")
                .should("contain", "/person/36/23");
            $row.should("contain", "bob smith");
            $row.should("contain", "123 Somewhere Road");
            $row.should("contain", "perfect");
            $row.should("contain", "LPA - PFA");
        });

        it("it cannot find any results", () => {
            cy.visit("/search", {
                qs: { 'term': 'abcde' }
            });
            cy.contains("No cases were found");
        });
    });

    describe("Search with filters", () => {
        beforeEach(() => {
            cy.visit("/search", {
                qs: { 'term': 'jack' }
            });
        });

        it("finds an attorney when filtered", () => {
            cy.contains("You searched for: jack");
            cy.contains("Showing 1 to 1 of 1 cases");
            cy.contains("Donor (1)");
            const $row = cy.get("table > tbody > tr");
            $row
                .contains("7000-8548-8461")
                .should("have.attr", "href")
                .should("contain", "/person/36/23");
            $row.should("contain", "jack smith");
            $row.should("contain", "123 Somewhere Road");
            $row.should("contain", "perfect");
            $row.should("contain", "LPA - PFA");
        });
    });

    describe("Search donor not associated with case", () => {
        beforeEach(() => {
            cy.visit("/search", {
                qs: { 'term': 'daniel' }
            });
        });

        it("finds a donor not associated with a case", () => {
            cy.contains("Showing 1 to 1 of 1 cases");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "Not associated with a case");
            $row
                .contains("Daniel Jones")
                .should("have.attr", "href")
                .should("contain", "/person/33");
            $row.should("contain", "22 Test Road");
        });
    });

    describe("Search deleted case", () => {
        beforeEach(() => {
            cy.visit("/search?term=700000005555");
        });

        it("finds a deleted case when searching by uid", () => {
            cy.contains("Search results");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "7000-0000-5555");
            $row.should("contain", "LPA");
            $row.should("contain", "02/12/2022");
            $row.should("contain", "return - unpaid");
            $row.should("contain", "LPA was not paid for after 12 months");
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
