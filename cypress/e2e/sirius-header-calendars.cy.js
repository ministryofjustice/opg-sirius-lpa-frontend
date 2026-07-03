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

    cy.addMock("/lpa-api/v1/persons/1/references", "GET", {
      status: 200,
      body: [
        {
          referenceId: 123,
        },
      ],
    });
  });

  it("displays the calendars panel", () => {
    const freezeDate = new Date(2025, 0, 10);
    cy.clock(freezeDate.getTime());
    cy.visit("/donor/1/documents");
    cy.get("#header-button-calendars").click();
    cy.get(".panel-calendar").should("be.visible");

    cy.get('[id^="calendar-month-"]').should("have.length", 3);
  });

  it("shows the working-days calculator and updates readonly fields by mode", () => {
    cy.visit("/donor/1/documents");
    cy.get("#header-button-calendars").click();

    cy.get(".panel-calendar").should("be.visible");
    cy.contains("h3", "Difference Calculator").should("be.visible");

    cy.get("#mode-enddate").should("be.checked");
    cy.get("#calc-enddate").should("have.attr", "readonly");
    cy.get("#calc-startdate").should("not.have.attr", "readonly");
    cy.get("#calc-numworkingdays").should("not.have.attr", "readonly");

    cy.get("#mode-startdate").click();
    cy.get("#mode-startdate").should("be.checked");
    cy.get("#calc-startdate").should("have.attr", "readonly");
    cy.get("#calc-enddate").should("not.have.attr", "readonly");
    cy.get("#calc-numworkingdays").should("not.have.attr", "readonly");

    cy.get("#mode-numworkingdays").click();
    cy.get("#mode-numworkingdays").should("be.checked");
    cy.get("#calc-numworkingdays").should("have.attr", "readonly");
    cy.get("#calc-startdate").should("not.have.attr", "readonly");
    cy.get("#calc-enddate").should("not.have.attr", "readonly");

    cy.get("#mode-enddate").click();
    cy.get("#mode-enddate").should("be.checked");
    cy.get("#calc-enddate").should("have.attr", "readonly");
    cy.get("#calc-startdate").should("not.have.attr", "readonly");
    cy.get("#calc-numworkingdays").should("not.have.attr", "readonly");
  });
});
