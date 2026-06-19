describe("Calendars on the header bar", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/dates/bank-holidays", "GET", {
      status: 200,
      body: {
        2025: {
          "New Year": "2025-01-01T00:00:00+00:00",
        },
      },
    });

    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {},
    });

    cy.addMock("/lpa-api/v1/persons/1/cases", "GET", {
      status: 200,
      body: {},
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {},
      },
    );

    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {},
    });
  });

  it("displays the calendars panel", () => {
    const freezeDate = new Date(2025, 0, 10);
    cy.clock(freezeDate.getTime());
    cy.visit("/donor/1/documents");
    cy.get("#header-button-calendars").click();
    cy.get(".panel-calendar").should("be.visible");
    cy.get('[data-bank-holiday="true"]');
    cy.get(":nth-child(1) > .calendar > .current-month > span").should(
      "contain.text",
      "December 2024",
    );
    cy.get(":nth-child(2) > .calendar > .current-month > span").should(
      "contain.text",
      "January 2025",
    );
    cy.get(":nth-child(3) > .calendar > .current-month > span").should(
      "contain.text",
      "February 2025",
    );
    cy.get(".prev-month-0").click();
    cy.get(":nth-child(1) > .calendar > .current-month > span").should(
      "contain.text",
      "November 2024",
    );
    cy.get(".next-month-1").click();
    cy.get(":nth-child(2) > .calendar > .current-month > span").should(
      "contain.text",
      "February 2025",
    );
  });
});
