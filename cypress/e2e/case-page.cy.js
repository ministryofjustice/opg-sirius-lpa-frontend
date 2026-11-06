describe("Case page", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {
        id: 1,
      },
    });

    cy.addMock("/lpa-api/v1/persons/1/cases", "GET", {
      status: 200,
      body: {
        cases: [
          {
            caseType: "LPA",
            caseSubtype: "pfa",
            id: 34,
            uId: "7000-1234-1234",
          },
        ],
      },
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:34&limit=999",
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

    cy.addMock(`/lpa-api/v1/lpas/34/draft-count`, "GET", {
      status: 200,
      body: {
        draftCount: 1,
      },
    });

    cy.addMock(
      `/lpa-api/v1/cases/34/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC`,
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/persons/1/references", "GET", {
      status: 200,
      body: [
        {
          referenceId: 123,
        },
      ],
    });

    cy.addMock("/config", "GET", {
      status: 200,
      body: {
        featureToggles: {
          poasNewCasePage: true,
        },
      },
    });

    cy.addMock("/lpa-api/v1/persons/1/events?&sort=id:desc&limit=999", "GET", {
      status: 200,
      body: {},
    });

    cy.visit("/donor/1?uid[]=7000-1234-1234");
  });

  it("should display the case page", () => {
    cy.get(".sirius-header").should("exist");
    cy.get(".action-panel__inner").should("exist");
    cy.get(".app-donor-details-table").should("exist");
    cy.get(".app-timeline-container").should("exist");
    cy.contains("Case Summary");
    cy.contains("Tasks");
  });
});
