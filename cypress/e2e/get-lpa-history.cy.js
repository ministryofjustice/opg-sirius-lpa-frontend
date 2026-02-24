describe("View LPA history timeline", () => {
  beforeEach(() => {
    const taskClass = String.raw`Opg\Core\Model\Entity\Task\Task`;
    const warningClass = String.raw`Opg\Core\Model\Entity\Warning\Warning`;
    const lpaClass = String.raw`Opg\Core\Model\Entity\CaseItem\PowerOfAttorney\Lpa`;
    const paymentClass = String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Payment\Payment`;
    const outgoingDocumentClass = String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Document\OutgoingDocument`;

    cy.addMock("/lpa-api/v1/persons/1/events?&sort=id:desc&limit=999", "GET", {
      status: 200,
      body: {
        limit: 999,
        metadata: {
          caseIds: [
            {
              id: 105,
              total: 4, //events on caseId
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
            {
              sourceType: "Payment",
              total: 1,
            },
          ],
        },
        pages: {
          current: 1,
          total: 1,
        },
        total: 4,
        events: [
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
              _class: warningClass,
              closedBy: [],
              warningText: "Test",
              warningType: "Complaint Received",
            },
            createdOn: "2026-01-22T16:23:29+00:00",
            hash: "N6R",
          },
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
              _class: taskClass,
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
              _class: lpaClass,
            },
            createdOn: "2026-01-15T04:10:55+00:00",
            hash: "JI7",
          },
          {
            id: 238,
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
            sourceType: "Payment",
            type: "DEL",
            changeSet: [],
            entity: {
              _class: paymentClass,
              amount: 2345,
              source: "CHEQUE",
              paymentDate: "2006-01-02T15:04:05+00:00",
            },
            createdOn: "2026-01-22T16:23:29+00:00",
            hash: "N7R",
          },
          {
            id: 239,
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
            sourceType: "OutgoingDocument",
            type: "INS",
            changeSet: [],
            entity: {
              _class: outgoingDocumentClass,
              friendlyDescription: "Joe Bloggs - Letter sent to donor",
            },
            sourceDocument: {
              UUID: "123e4567-e89b-12d3-a456-426614174000",
            },
            createdOn: "2026-01-23T16:23:29+00:00",
            hash: "N8R",
          }
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
          total: 4,
          events: [
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
                _class: warningClass,
                closedBy: [],
                warningText: "Test",
                warningType: "Complaint Received",
              },
              createdOn: "2026-01-22T16:23:29+00:00",
              hash: "N6R",
            },
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
                _class: taskClass,
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
                _class: lpaClass,
              },
              createdOn: "2026-01-15T04:10:55+00:00",
              hash: "JI7",
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
                _class: warningClass,
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

    cy.addMock(
      "/lpa-api/v1/persons/1/events?filter=case:105,case:106,case:107,sourceType:Warning&sort=id:asc&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            caseIds: [
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
                sourceType: "Warning",
                total: 2,
              },
            ],
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 2,
          events: [
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
                _class: warningClass,
                closedBy: [],
                warningType: "Payment Required",
              },
              createdOn: "2026-01-01T09:35:18+00:00",
              hash: "GQN",
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
                _class: warningClass,
                closedBy: [],
                warningText: "Test",
                warningType: "Complaint Received",
              },
              createdOn: "2026-01-22T16:23:29+00:00",
              hash: "N6R",
            },
          ],
        },
      },
    );

    cy.visit("/donor/1/history");
  });

  it("can view variable header content", () => {
    cy.get(".moj-timeline__item")
      .should("have.length", 5)
      .then(($items) => {
        expect($items.eq(0)).to.contain.text("Warning");
        expect($items.eq(1)).to.contain.text("Task");
        expect($items.eq(2)).to.contain.text("LPA (Create / Edit)");
        expect($items.eq(3)).to.contain.text("Payment");
        expect($items.eq(4)).to.contain.text("Outbound document");
        cy.wrap($items.eq(1)).find(".moj-alert--warning").should("not.exist");
        cy.wrap($items.eq(0)).find(".moj-alert--warning").should("exist");
      });
  });

  it("can view variable footer content", () => {
    cy.get(".moj-timeline__item")
      .should("have.length", 5)
      .then(($items) => {
        const normalise = (el) =>
          Cypress.$(el).text().replaceAll(/\s+/g, " ").trim();

        expect(normalise($items[0])).to.include(
          "Deleted by Team Less – 03001234567",
        );
        expect(normalise($items[1])).to.include(
          "Created by Marty Test (Casework Team) – 0300 300 0300",
        );
        expect(normalise($items[2])).to.include("Updated by deleted user");
      });
  });

  it("can view variable case details content", () => {
    cy.visit("/donor/1/history?id[]=105&id[]=106&id[]=107");

    cy.get(".moj-timeline__item")
      .should("have.length", 4)
      .then(($items) => {
        cy.wrap($items.eq(0))
          .should("contain.text", "EPA 7000-7000-7000")
          .find(".colour-govuk-brown")
          .should("exist");

        cy.wrap($items.eq(1))
          .should("contain.text", "PFA 7000-9000-7000")
          .find(".colour-govuk-turquoise")
          .should("exist");

        cy.wrap($items.eq(2))
          .should("contain.text", "HW 7000-9000-6000")
          .find(".colour-govuk-grass-green")
          .should("exist");

        cy.wrap($items.eq(3))
          .should("not.contain.text", "EPA")
          .and("not.contain.text", "LPA")
          .and("not.contain.text", "HW");
      });
  });

  it("can filter", () => {
    cy.visit("/donor/1/history?id[]=105&id[]=106&id[]=107");

    cy.contains("(showing all 4 items)");
    cy.contains("Ascending").click();
    cy.contains("Warning (2)").click();
    cy.contains("Apply filters").click();

    cy.contains("(showing 2 of 4 items)");
    cy.get(".moj-timeline__item")
      .should("have.length", 2)
      .then(($items) => {
        expect($items.eq(0)).to.contain.text("Warning");
        expect($items.eq(1)).to.contain.text("Warning");

        cy.wrap($items.eq(0))
          .should("not.contain.text", "EPA")
          .and("not.contain.text", "LPA")
          .and("not.contain.text", "HW");

        cy.wrap($items.eq(1))
          .should("contain.text", "EPA 7000-7000-7000")
          .find(".colour-govuk-brown")
          .should("exist");
      });
  });

  it("can view payment deleted event", () => {
    cy.get(".moj-timeline__item")
      .eq(3)
      .should("contain.text", "Deleted - £23.45 paid by cheque on 02/01/2006");
  });

  it("can view outbound document event", () => {
    cy.get(".moj-timeline__item")
      .eq(4)
      .should("contain.text", "Joe Bloggs - Letter sent to donor")
      .find("a")
      .should("have.attr", "href", "/lpa#/donor/1/documents?docUuid=123e4567-e89b-12d3-a456-426614174000");
  });
});
