describe("Search", () => {
    describe("Searching by name", () => {

        it("finds a person not associated with a case", () => {
            cy.visit("/search?term=bob");
            cy.contains("You searched for: bob");
            cy.contains("Showing 1 to 1 of 1 results");
            cy.contains("Donor (1)");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "Bob Smith");
            $row.should("contain", "17/03/1990");
            $row.should("contain", "123 Somewhere Road");
            $row.should("contain", "perfect");
            $row.should("contain", "LPA - PFA");
            $row
                .contains("7000-5382-4438")
                .should("have.attr", "href")
                .should("contain", "/person/36/23");
        });

        it("finds a person with more than one case", () => {
            cy.visit("/search?term=harry");
            cy.contains("You searched for: harry");
            cy.contains("Showing 1 to 1 of 1 results");
            cy.contains("Donor (1)");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "Harry Jones");
            $row.should("contain", "123 Somewhere Road");
            $row.should("contain", "perfect");
            $row.should("contain", "pending");
            $row.should("contain", "LPA - PFA");
            $row.should("contain", "LPA - HW");
            cy.contains("7000-5382-4438")
                .should("have.attr", "href")
                .should("contain", "/person/36/23");
            cy.contains("7000-5382-8764")
                .should("have.attr", "href")
                .should("contain", "/person/36/24");
        });

        it("finds a deputy", () => {
            cy.visit("/search?term=fred");
            cy.contains("You searched for: fred");
            cy.contains("Showing 1 to 1 of 1 results");
            cy.contains("Deputy (1)");
            const $row = cy.get("table > tbody > tr");
            $row.should("contain", "Fred Jones");
            $row.should("contain", "Deputy");
            $row.should("contain", "100 Random Road");
            $row.should("contain", "ORDER - HW");
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

        it("can apply and remove filters", () => {
            cy.contains(".govuk-checkboxes__item", "Attorney").find("input").check();
            cy.contains(".govuk-checkboxes__item", "Trust corporation").find("input").check();
            cy.get("button[type=submit]").click();
            cy.contains(".moj-filter__tag", "Attorney");
            cy.contains(".moj-filter__tag", "Trust Corporation");
            cy.contains(".moj-filter__selected-heading", "Clear filters").find("a").click();
            cy.get('.moj-filter__tag').should('not.exist');
        });
    });

    describe("Search deleted case", () => {
       beforeEach(() => {
           cy.visit("/search?term=700000005555");
       });

       it("finds a deleted case", () => {
           const $row = cy.get("table > tbody > tr");
           $row.should("contain", "LPA was not paid for after 12 months");
           $row.should("contain", "return - unpaid");
           $row.should("contain", "A987654321");
           $row.should("contain", "02/12/2022");
           $row.should("contain", "LPA");
           $row.should("contain", "deleted");
           $row.should("contain", "7000-0000-5555");
       });
    });
});
