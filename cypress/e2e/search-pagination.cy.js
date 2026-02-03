describe("Search pagination", () => {
  describe("Next and Previous links", () => {
    beforeEach(() => {
      // We have to mock a search with enough results to trigger the pagination (currently > 25)
      const mockResults = Array.from({ length: 80 }, (_, i) => ({
        id: i + 1,
        uId: `7000-0000-${String(i + 1).padStart(4, "0")}`,
        personType: i % 3 === 0 ? "Donor" : i % 3 === 1 ? "Attorney" : "Deputy",
        firstname: "Test",
        surname: `Person${i + 1}`,
        dob: "01/01/1980",
        addressLine1: `${i + 1} Test Road`,
        cases: [
          {
            id: i + 100,
            uId: `7000-1000-${String(i + 1).padStart(4, "0")}`,
            caseType: "LPA",
            caseSubtype: "pfa",
            status: "Pending",
          },
        ],
      }));

      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: {
            personType: {
              Donor: 27,
              Attorney: 27,
              Deputy: 26,
            },
          },
          results: mockResults,
          total: {
            count: 80,
          },
        },
      });
    });

    it("displays pagination controls when results exceed page size", () => {
      cy.visit("/search?term=test");

      cy.contains("Showing 1 to 25");
      cy.get(".govuk-pagination").should("exist");
      cy.get(".govuk-pagination__next").should("exist");
      cy.get(".govuk-pagination__prev").should("not.exist"); // Not on page 1
    });

    it("next link has correct href and navigates to page 2", () => {
      cy.visit("/search?term=test");

      cy.get(".govuk-pagination__next a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=2");

      cy.get(".govuk-pagination__next a").click();

      cy.url().should("include", "page=2");
      cy.contains("Showing 26 to 50");
    });

    it("previous link has correct href and navigates to page 1", () => {
      cy.visit("/search?term=test&page=2");

      cy.contains("Showing 26 to 50");

      cy.get(".govuk-pagination__prev a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=1");

      cy.get(".govuk-pagination__prev a").click();

      cy.url().should("not.include", "page=2");
      cy.contains("Showing 1 to 25");
    });

    it("next and previous links work correctly on middle pages", () => {
      cy.visit("/search?term=test&page=2");

      cy.get(".govuk-pagination__next").should("exist");
      cy.get(".govuk-pagination__prev").should("exist");

      cy.get(".govuk-pagination__next a")
        .should("have.attr", "href")
        .and("include", "page=3");

      cy.get(".govuk-pagination__prev a")
        .should("have.attr", "href")
        .and("include", "page=1");
    });

    it("hides next link on last page", () => {
      cy.visit("/search?term=test&page=4");

      cy.contains("Showing 76 to 80 of 80 results");

      cy.get(".govuk-pagination__next").should("not.exist");

      cy.get(".govuk-pagination__prev").should("exist");
      cy.get(".govuk-pagination__prev a")
        .should("have.attr", "href")
        .and("include", "page=3");
    });

    it("preserves filters in next/previous links", () => {
      cy.visit("/search?term=test");

      cy.contains("label", "Donor").click();
      cy.contains(".moj-filter__tag", "Donor");

      cy.get(".govuk-pagination__next a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "person-type=Donor")
        .and("include", "page=2");

      cy.get(".govuk-pagination__next a").click();

      cy.url().should("include", "person-type=Donor");
      cy.contains(".moj-filter__tag", "Donor");

      cy.get(".govuk-pagination__prev a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "person-type=Donor")
        .and("include", "page=1");
    });

    it("numbered page links work alongside next/previous", () => {
      cy.visit("/search?term=test&page=2");

      cy.get(".govuk-pagination__item a").contains("3").click();

      cy.url().should("include", "page=3");
      cy.contains("Showing 51 to 75");

      cy.get(".govuk-pagination__next").should("exist");
      cy.get(".govuk-pagination__prev").should("exist");
    });
  });

  describe("Pagination with no results", () => {
    beforeEach(() => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: {},
          results: [],
          total: {
            count: 0,
          },
        },
      });

      cy.visit("/search?term=nonexistent");
    });

    it("does not show pagination controls", () => {
      cy.contains("Showing 1 to 0");
      cy.get(".govuk-pagination").should("not.exist");
    });
  });

  describe("Pagination with single page", () => {
    beforeEach(() => {
      const mockResults = Array.from({ length: 10 }, (_, i) => ({
        id: i + 1,
        uId: `7000-0000-${String(i + 1).padStart(4, "0")}`,
        personType: "Donor",
        firstname: "Test",
        surname: `Person${i + 1}`,
        dob: "01/01/1980",
        addressLine1: `${i + 1} Test Road`,
        cases: [
          {
            id: i + 100,
            uId: `7000-1000-${String(i + 1).padStart(4, "0")}`,
            caseType: "LPA",
            caseSubtype: "pfa",
            status: "Pending",
          },
        ],
      }));

      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: { personType: { Donor: 10 } },
          results: mockResults,
          total: {
            count: 10,
          },
        },
      });

      cy.visit("/search?term=test");
    });

    it("does not show pagination controls with results under page size", () => {
      cy.contains("Showing 1 to 10 of 10 results");
      cy.get(".govuk-pagination__next").should("not.exist");
      cy.get(".govuk-pagination__prev").should("not.exist");
      cy.get(".govuk-pagination__list").should("not.exist");
    });
  });
});
