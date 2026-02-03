describe("Search pagination - Next and Previous link hrefs", () => {
  describe("With paginated results", () => {
    beforeEach(() => {
      // Create a mock with enough results to trigger pagination
      const generateResults = (count) => {
        return Array.from({ length: count }, (_, i) => ({
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
      };

      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: { personType: { Donor: 80 } },
          results: generateResults(80),
          total: { count: 80 },
        },
      });
    });

    it("next link on page 1 has correct href with page=2", () => {
      cy.visit("/search?term=test");

      cy.get(".govuk-pagination__next a")
        .should("exist")
        .and("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=2")
        .and("not.equal", "#");
    });

    it("previous link on page 2 has correct href with page=1", () => {
      cy.visit("/search?term=test&page=2");

      cy.get(".govuk-pagination__prev a")
        .should("exist")
        .and("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=1")
        .and("not.equal", "#");
    });

    it("next link on page 2 has correct href with page=3", () => {
      cy.visit("/search?term=test&page=2");

      cy.get(".govuk-pagination__next a")
        .should("exist")
        .and("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=3")
        .and("not.equal", "#");
    });

    it("preserves search term and filters in pagination hrefs", () => {
      cy.visit("/search?term=test&person-type=Donor");

      cy.get(".govuk-pagination__next a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "person-type=Donor")
        .and("include", "page=2");
    });

    it("numbered page links have same format as next/previous links", () => {
      cy.visit("/search?term=test");

      cy.get(".govuk-pagination__item a")
        .first()
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=");

      cy.get(".govuk-pagination__next a")
        .should("have.attr", "href")
        .and("include", "term=test")
        .and("include", "page=2");
    });
  });

  describe("Edge cases", () => {
    it("no next link on last page", () => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: { personType: { Donor: 30 } },
          results: Array.from({ length: 30 }, (_, i) => ({
            id: i + 1,
            uId: `7000-${i + 1}`,
            personType: "Donor",
            firstname: "Test",
            surname: `Person${i}`,
            cases: [],
          })),
          total: { count: 30 },
        },
      });

      cy.visit("/search?term=test&page=2");

      cy.get(".govuk-pagination__next").should("not.exist");

      cy.get(".govuk-pagination__prev a")
        .should("exist")
        .and("have.attr", "href")
        .and("not.equal", "#");
    });

    it("no previous link on first page", () => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: { personType: { Donor: 30 } },
          results: Array.from({ length: 30 }, (_, i) => ({
            id: i + 1,
            uId: `7000-${i + 1}`,
            personType: "Donor",
            firstname: "Test",
            surname: `Person${i}`,
            cases: [],
          })),
          total: { count: 30 },
        },
      });

      cy.visit("/search?term=test");

      cy.get(".govuk-pagination__prev").should("not.exist");

      cy.get(".govuk-pagination__next a")
        .should("exist")
        .and("have.attr", "href")
        .and("not.equal", "#");
    });

    it("no pagination controls with single page of results", () => {
      cy.addMock("/lpa-api/v1/search/persons", "POST", {
        status: 200,
        body: {
          aggregations: { personType: { Donor: 10 } },
          results: Array.from({ length: 10 }, (_, i) => ({
            id: i + 1,
            uId: `7000-${i + 1}`,
            personType: "Donor",
            firstname: "Test",
            surname: `Person${i}`,
            cases: [],
          })),
          total: { count: 10 },
        },
      });

      cy.visit("/search?term=test");

      // With only 10 results (< 25), no pagination should render
      cy.get(".govuk-pagination__next").should("not.exist");
      cy.get(".govuk-pagination__prev").should("not.exist");
    });
  });
});
