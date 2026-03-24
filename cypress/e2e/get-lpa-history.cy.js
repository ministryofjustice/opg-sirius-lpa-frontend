describe("View LPA history timeline", () => {
  beforeEach(() => {
    const taskClass = String.raw`Opg\Core\Model\Entity\Task\Task`;
    const warningClass = String.raw`Opg\Core\Model\Entity\Warning\Warning`;
    const lpaClass = String.raw`Opg\Core\Model\Entity\CaseItem\PowerOfAttorney\Lpa`;

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
      .should("have.length", 3)
      .then(($items) => {
        expect($items.eq(0)).to.contain.text("Warning");
        expect($items.eq(1)).to.contain.text("Task");
        expect($items.eq(2)).to.contain.text("LPA (Create / Edit)");
        cy.wrap($items.eq(1)).find(".moj-alert--warning").should("not.exist");
        cy.wrap($items.eq(0)).find(".moj-alert--warning").should("exist");
      });
  });

  it("can view variable footer content", () => {
    cy.get(".moj-timeline__item")
      .should("have.length", 3)
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

  it("can view phone number event content", () => {
    cy.addMock("/lpa-api/v1/persons/2/events?&sort=id:desc&limit=999", "GET", {
      status: 200,
      body: {
        limit: 999,
        metadata: {
          caseIds: [
            {
              id: 900,
              total: 1,
            },
            {
              total: 1,
            },
          ],
          sourceTypes: [
            {
              sourceType: "PhoneNumber",
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
            id: 170,
            owningCase: {
              id: 900,
              uId: "7000-0000-0009",
              caseSubtype: "pfa",
              caseType: "LPA",
            },
            user: {
              id: 5,
              phoneNumber: "030030000300",
              teams: [],
              displayName: "OPG User",
              deleted: false,
              email: "opg@test.gov.uk",
            },
            sourceType: "PhoneNumber",
            sourcePhoneNumber: {
              phoneNumber: "12345 678910",
              type: "Work",
            },
            type: "UPD",
            changeSet: {
              phoneNumber: ["12345", "12345 678910"],
            },
            createdOn: "2026-03-06T14:39:20+00:00",
            hash: "ABC",
          },
          {
            id: 499,
            owningCase: {
              id: 900,
              uId: "7000-0000-0009",
              caseSubtype: "pfa",
              caseType: "LPA",
            },
            user: {
              id: 5,
              phoneNumber: "030030000300",
              teams: [],
              displayName: "OPG User",
              deleted: false,
              email: "opg@test.gov.uk",
            },
            sourceType: "PhoneNumber",
            sourcePhoneNumber: {
              phoneNumber: "12345",
              type: "Work",
            },
            type: "INS",
            changeSet: [],
            createdOn: "2026-01-22T10:30:01+00:00",
            hash: "AB",
          },
        ],
      },
    });

    cy.visit("/donor/2/history");

    cy.get(".moj-timeline__item")
      .should("have.length", 2)
      .then(($items) => {
        const normalise = (el) =>
          Cypress.$(el).text().replaceAll(/\s+/g, " ").trim();
        expect(normalise($items[0])).to.include(
          "Phone number changed from 12345 to 12345 678910",
        );
        expect(normalise($items[1])).to.include(
          "Phone number changed to 12345",
        );
      });
  });

  it("can view donor event content", () => {
    cy.addMock("/lpa-api/v1/persons/2/events?&sort=id:desc&limit=999", "GET", {
      status: 200,
      body: {
        limit: 999,
        metadata: {
          caseIds: [
            {
              total: 2,
            },
          ],
          sourceTypes: [
            {
              sourceType: "Donor",
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
            id: 176,
            user: {
              id: 5,
              phoneNumber: "030030000300",
              teams: [],
              displayName: "OPG User",
              deleted: false,
              email: "opg@test.gov.uk",
            },
            sourceType: "Donor",
            sourcePerson: {
              id: 17,
              uId: "7000-0000-0009",
              firstname: "Test",
              surname: "Case",
            },
            type: "UPD",
            changeSet: {
              firstname: ["Testing", "Test"],
              surname: ["Casing", "Case"],
              email: ["test@test.com", "test@testcase.com"],
              salutation: ["Mr", "Mrs"],
              correspondenceByEmail: [false, true],
              dob: {
                1: {
                  date: "1999-05-09 00:00:00.000000",
                  timezone_type: 3,
                  timezone: "UTC",
                },
              },
            },
            createdOn: "2026-01-31T14:39:20+00:00",
            hash: "AAA",
          },
          {
            id: 175,
            user: {
              id: 5,
              phoneNumber: "030030000300",
              teams: [],
              displayName: "OPG User",
              deleted: false,
              email: "opg@test.gov.uk",
            },
            sourceType: "Donor",
            sourcePerson: {
              id: 17,
              uId: "7000-0000-0009",
              firstname: "Test",
              surname: "Case",
            },
            type: "INS",
            changeSet: [],
            entity: {
              _class: String.raw`Opg\Core\Model\Entity\CaseActor\Donor`,
              email: "test@test.com",
              firstname: "Testing",
              id: 17,
              salutation: "Mr",
              surname: "Casing",
              uId: 700000000009,
            },
            createdOn: "2026-01-22T10:30:01+00:00",
            hash: "AZ",
          },
        ],
      },
    });

    cy.visit("/donor/2/history");

    cy.get(".moj-timeline__item")
      .first()
      .within(() => {
        cy.contains("Salutation: Mr changed to: Mrs");
        cy.contains("First name: Testing changed to: Test");
        cy.contains("Surname: Casing changed to: Case");
        cy.contains("Date of birth: 09/05/1999");
        cy.contains("Email: test@test.com changed to: test@testcase.com");
        cy.contains("Correspondence by email: false changed to: true");
      });

    cy.get(".moj-timeline__item")
      .last()
      .within(() => {
        cy.contains("Testing Casing");
      });
  });

  it("can filter", () => {
    cy.visit("/donor/1/history?id[]=105&id[]=106&id[]=107");

    cy.contains("Apply filters").should("not.be.visible");

    cy.contains("(showing all 4 items)");
    cy.contains("Ascending").click();
    cy.contains("Warning (2)").click();

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

  describe("Filter panel visibility", () => {
    it("shows the filter panel by default", () => {
      cy.get(".moj-filter-layout__filter").should("be.visible");
      cy.get("div[data-filter-summary]").should("not.be.visible");
    });

    it("hides the filter panel when Hide filters is clicked", () => {
      cy.contains(".govuk-button", "Hide filters").click();
      cy.get(".moj-filter-layout__filter").should("not.be.visible");
      cy.get("div[data-filter-summary]").should("be.visible");
    });

    it("shows the filter panel again when Show filters is clicked", () => {
      cy.contains(".govuk-button", "Hide filters").click();
      cy.contains(".govuk-button", "Show filters").click();
      cy.get(".moj-filter-layout__filter").should("be.visible");
      cy.get("div[data-filter-summary]").should("not.be.visible");
    });
  });
});
