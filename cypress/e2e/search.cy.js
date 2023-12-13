describe("Search", () => {
  describe("Searching by name", () => {
    beforeEach(() => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: {
            personType: {
              Donor: 2,
              Deputy: 1,
            },
          },
          results: [
            {
              id: 36,
              uId: "7000-8548-8461",
              personType: "Donor",
              firstname: "Bob",
              surname: "Smith",
              dob: "17/03/1990",
              addressLine1: "123 Somewhere Road",
              cases: [
                {
                  id: 23,
                  uId: "7000-5382-4438",
                  caseType: "LPA",
                  caseSubtype: "pfa",
                  status: "Perfect",
                },
              ],
            },
            {
              id: 36,
              uId: "7000-8548-8461",
              personType: "Donor",
              firstname: "Bob",
              surname: "Jones",
              addressLine1: "123 Somewhere Road",
              cases: [
                {
                  id: 23,
                  uId: "7000-5382-4438",
                  caseType: "LPA",
                  caseSubtype: "pfa",
                  status: "Perfect",
                },
                {
                  id: 24,
                  uId: "7000-5382-8764",
                  caseType: "LPA",
                  caseSubtype: "hw",
                  status: "Pending",
                },
              ],
            },
            {
              id: 65,
              uId: "7000-6509-8813",
              personType: "Deputy",
              firstname: "Bob",
              surname: "Rogers",
              addressLine1: "100 Random Road",
              cases: [
                {
                  id: 48,
                  uId: "7000-5113-1871",
                  caseType: "ORDER",
                  caseSubtype: "hw",
                },
              ],
            },
          ],
          total: {
            count: 3,
          },
        },
      });

      cy.visit("/search?term=bob");

      cy.contains("You searched for: bob");
      cy.contains("Showing 1 to 3 of 3 results");
      cy.contains("Donor (2)");
      cy.contains("Deputy (1)");
    });

    it("finds a person not associated with a case", () => {
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
      const $row = cy.get("table > tbody > tr");
      $row.should("contain", "Bob Jones");
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
      const $row = cy.get("table > tbody > tr");
      $row.should("contain", "Bob Rogers");
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
      cy.get(".moj-filter-layout__filter").should("not.be.visible");
      cy.contains(".govuk-button", "Show filters").click();
      cy.get(".moj-filter-layout__filter").should("be.visible");
    });

    it("can apply and remove filters", () => {
      // Checks the button is hidden because js is enabled
      cy.contains("Apply filters").should("not.be.visible");
      cy.contains(".govuk-checkboxes__item", "Attorney").find("input").check();
      cy.contains(".govuk-checkboxes__item", "Trust corporation")
        .find("input")
        .check();
      cy.contains(".moj-filter__tag", "Attorney");
      cy.contains(".moj-filter__tag", "Trust Corporation");
      cy.contains(".moj-filter__selected-heading", "Clear filters")
        .find("a")
        .click();
      cy.get(".moj-filter__tag").should("not.exist");
    });
  });

  describe("Search deleted case", () => {
    beforeEach(() => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          total: {
            count: 0,
          },
        },
      });

      cy.visit("/search?term=700000005555");
    });

    it("finds a deleted case", () => {
      const $row = cy.get("table > tbody > tr");
      $row.should("contain", "LPA was not paid for after 12 months");
      $row.should("contain", "return - unpaid");
      $row.should("contain", "02/12/2022");
      $row.should("contain", "LPA");
      $row.should("contain", "deleted");
      $row.should("contain", "7000-0000-5555");
    });

    it("finds no cases (dropdown appears)", () => {
      cy.visit("/search?term=abcdefg");
      cy.get("#f-search-form-below-phase-banner").type("abcdefg").submit();
      const $row = cy.get(".sirius-search__item");
      $row.should("contain", "No cases were found");
    });
  });
});
