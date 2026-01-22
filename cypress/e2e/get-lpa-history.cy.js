describe("View LPA history timeline", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1/events?&sort=id:desc&limit=999", "GET", {
      status: 200,
      body: {
        limit: 999,
        metadata: {
          caseIds: [
            {
              id: 105,
              total: 3, //events on caseId
            },
            {
              total: 0, //events on person
            },
          ],
          sourceTypes: [
            {
              sourceType: "Lpa",
              total: 1,
            },
            {
              sourceType: "Task",
              total: 1,
            },
            {
              sourceType: "Warning",
              total: 1,
            },
          ],
        },
        pages: {
          current: 1,
          total: 1,
        },
        total: 3,
        events: [
          {
            id: 117,
            owningCase: {
              id: 105,
              uId: "7000-9000-7000",
              caseSubtype: "pfa",
              caseType: "LPA",
            },
            user: {
              id: 618,
              phoneNumber: "0300 300 0300",
              teams: [
                {
                  id: 71,
                  phoneNumber: "0300 3000300",
                  teams: [],
                  displayName: "Casework Team",
                  deleted: false,
                  email: "test@test.gov.uk",
                },
              ],
              displayName: "Marty Test",
              deleted: false,
              email: "marty.test@publicguardian.uk",
            },
            sourceType: "Task",
            type: "INS",
            changeSet: [],
            entity: {
              _class: "Opg\\Core\\Model\\Entity\\Task\\Task",
              assignee: {
                displayName: "Registration Team",
              },
              name: "Autogenerate letters and register LPA",
            },
            createdOn: "2026-01-16T04:10:55+00:00",
            hash: "JIG",
          },
          {
            id: 119,
            owningCase: {
              id: 105,
              uId: "7000-9000-7000",
              caseSubtype: "pfa",
              caseType: "LPA",
            },
            user: {
              id: 618,
              phoneNumber: "0300 300 0300",
              teams: [
                {
                  id: 70,
                  phoneNumber: "0300 3000300",
                  teams: [],
                  displayName: "Casework Team",
                  deleted: false,
                  email: "test@test.uk",
                },
              ],
              displayName: "Test Deleted",
              deleted: true,
              email: "deleted@publicguardian.gov.uk",
            },
            sourceType: "Lpa",
            type: "UPD",
            changeSet: {
              noticeGivenDate: {
                1: {
                  date: "2024-04-15 04:10:55.617326",
                  timezone_type: 3,
                  timezone: "UTC",
                },
              },
            },
            entity: {
              _class:
                "Opg\\Core\\Model\\Entity\\CaseItem\\PowerOfAttorney\\Lpa",
            },
            createdOn: "2026-01-15T04:10:55+00:00",
            hash: "JI7",
          },
          {
            id: 144,
            owningCase: {
              id: 105,
              uId: "7000-9000-7000",
              caseSubtype: "pfa",
              caseType: "LPA",
            },
            user: {
              id: 6619,
              phoneNumber: "03001234567",
              teams: [],
              displayName: "Team Less",
              deleted: false,
              email: "teamless@test.uk",
            },
            sourceType: "Warning",
            sourceWarning: {
              id: 707,
            },
            type: "DEL",
            changeSet: [],
            entity: {
              _class: "Opg\\Core\\Model\\Entity\\Warning\\Warning",
              closedBy: [],
              warningText: "Test",
              warningType: "Complaint Received",
            },
            createdOn: "2026-01-22T16:23:29+00:00",
            hash: "N6R",
          },
        ],
      },
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/events?filter=case:105,case:106,case:107&sort=id:desc&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            caseIds: [
              {
                id: 105,
                total: 1,
              },
              {
                id: 106,
                total: 1,
              },
              {
                id: 107,
                total: 1,
              },
              {
                total: 1,
              },
            ],
            sourceTypes: [
              {
                sourceType: "Lpa",
                total: 1,
              },
              {
                sourceType: "Task",
                total: 1,
              },
              {
                sourceType: "Warning",
                total: 2,
              },
            ],
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 3,
          events: [
            {
              id: 117,
              owningCase: {
                id: 105,
                uId: "7000-9000-7000",
                caseSubtype: "pfa",
                caseType: "LPA",
              },
              user: {
                id: 618,
                phoneNumber: "0300 300 0300",
                teams: [
                  {
                    id: 71,
                    phoneNumber: "0300 3000300",
                    teams: [],
                    displayName: "Casework Team",
                    deleted: false,
                    email: "test@test.gov.uk",
                  },
                ],
                displayName: "Marty Test",
                deleted: false,
                email: "marty.test@publicguardian.uk",
              },
              sourceType: "Task",
              type: "INS",
              changeSet: [],
              entity: {
                _class: "Opg\\Core\\Model\\Entity\\Task\\Task",
                assignee: {
                  displayName: "Registration Team",
                },
                name: "Autogenerate letters and register LPA",
              },
              createdOn: "2026-01-16T04:10:55+00:00",
              hash: "JIG",
            },
            {
              id: 119,
              owningCase: {
                id: 106,
                uId: "7000-9000-6000",
                caseSubtype: "hw",
                caseType: "LPA",
              },
              user: {
                id: 618,
                phoneNumber: "0300 300 0300",
                teams: [
                  {
                    id: 70,
                    phoneNumber: "0300 3000300",
                    teams: [],
                    displayName: "Casework Team",
                    deleted: false,
                    email: "test@test.uk",
                  },
                ],
                displayName: "Test Deleted",
                deleted: true,
                email: "deleted@publicguardian.gov.uk",
              },
              sourceType: "Lpa",
              type: "UPD",
              changeSet: {
                noticeGivenDate: {
                  1: {
                    date: "2024-04-15 04:10:55.617326",
                    timezone_type: 3,
                    timezone: "UTC",
                  },
                },
              },
              entity: {
                _class:
                  "Opg\\Core\\Model\\Entity\\CaseItem\\PowerOfAttorney\\Lpa",
              },
              createdOn: "2026-01-15T04:10:55+00:00",
              hash: "JI7",
            },
            {
              id: 144,
              owningCase: {
                id: 107,
                uId: "7000-7000-7000",
                caseSubtype: "pfa",
                caseType: "EPA",
              },
              user: {
                id: 6619,
                phoneNumber: "03001234567",
                teams: [],
                displayName: "Team Less",
                deleted: false,
                email: "teamless@test.uk",
              },
              sourceType: "Warning",
              sourceWarning: {
                id: 707,
              },
              type: "DEL",
              changeSet: [],
              entity: {
                _class: "Opg\\Core\\Model\\Entity\\Warning\\Warning",
                closedBy: [],
                warningText: "Test",
                warningType: "Complaint Received",
              },
              createdOn: "2026-01-22T16:23:29+00:00",
              hash: "N6R",
            },
            {
              id: 101,
              user: {
                id: 6619,
                phoneNumber: "03001234567",
                teams: [],
                displayName: "Team Less",
                deleted: false,
                email: "teamless@test.uk",
              },
              sourceType: "Warning",
              sourceWarning: {
                id: 503,
              },
              type: "INS",
              changeSet: [],
              entity: {
                _class: "Opg\\Core\\Model\\Entity\\Warning\\Warning",
                closedBy: [],
                warningType: "Payment Required",
              },
              createdOn: "2026-01-01T09:35:18+00:00",
              hash: "GQN",
            },
          ],
        },
      },
    );

    cy.visit("/donor/1/history");
  });

  it("can view variable header content", () => {
    cy.get(".moj-timeline__item")
      .should("have.length", 3)
      .then(($items) => {
        expect($items.eq(0)).to.contain.text("Task");
        expect($items.eq(1)).to.contain.text("LPA (Create / Edit)");
        expect($items.eq(2)).to.contain.text("Warning");
        cy.wrap($items.eq(2)).find(".moj-alert--warning").should("exist");
      });
  });

  it("can view variable footer content", () => {
    cy.get(".moj-timeline__item")
      .should("have.length", 3)
      .then(($items) => {
        const normalise = (el) =>
          Cypress.$(el).text().replace(/\s+/g, " ").trim();

        expect(normalise($items[0])).to.include(
          "Created by Marty Test (Casework Team) – 0300 300 0300",
        );

        expect(normalise($items[1])).to.include("Updated by deleted user");

        expect(normalise($items[2])).to.include(
          "Deleted by Team Less – 03001234567",
        );
      });
  });

  it("can view variable case details content", () => {
    cy.visit("/donor/1/history?id[]=105&id[]=106&id[]=107");

    cy.get(".moj-timeline__item")
      .should("have.length", 4)
      .then(($items) => {
        cy.wrap($items.eq(0))
          .should("contain.text", "PFA 7000-9000-7000")
          .find(".colour-govuk-turquoise")
          .should("exist");

        cy.wrap($items.eq(1))
          .should("contain.text", "HW 7000-9000-6000")
          .find(".colour-govuk-grass-green")
          .should("exist");

        cy.wrap($items.eq(2))
          .should("contain.text", "EPA 7000-7000-7000")
          .find(".colour-govuk-brown")
          .should("exist");

        cy.wrap($items.eq(3))
          .should("not.contain.text", "EPA")
          .and("not.contain.text", "LPA")
          .and("not.contain.text", "HW");
      });
  });
});
