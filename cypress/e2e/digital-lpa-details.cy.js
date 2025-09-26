import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/333", "GET", {
      status: 200,
      body: {
        id: 333,
        uId: "M-DIGI-LPA3-3333",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Processing",
        expectedPaymentTotal: 8200,
      },
    });

    cy.addMock("/lpa-api/v1/tasks/1", "GET", {
      status: 200,
      body: {
        caseItems: [
          {
            caseType: "DIGITAL_LPA",
            uId: "M-DIGI-LPA3-3333",
          },
        ],
        dueDate: "10/01/2022",
        id: 1,
        name: "Create physical case file",
        status: "Not Started",
      },
    });

    const mocks = Promise.allSettled([
      digitalLpas.get("M-DIGI-LPA3-3333", {
        "opg.poas.sirius": {
          id: 333,
          donor: {
            id: 33,
          },
          linkedDigitalLpas: [
            {
              uId: "M-DIGI-LPA3-3334",
              caseSubtype: "personal-welfare",
              status: "Draft",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-DIGI-LPA3-3335",
              caseSubtype: "personal-welfare",
              status: "Registered",
              createdDate: "02/11/2023",
            },
          ],
        },
        "opg.poas.lpastore": {
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
              appointmentType: "original",
              cannotMakeJointDecisions: true,
            },
            {
              firstNames: "Volo",
              lastName: "McSpolo",
              status: "active",
              appointmentType: "original",
              signedAt: "2022-12-20T12:02:43Z",
            },
            {
              firstNames: "Susanna",
              lastName: "Kaysen",
              status: "removed",
              appointmentType: "original",
            },
            {
              firstNames: "Philomena",
              lastName: "Guinea",
              status: "inactive",
              appointmentType: "replacement",
            },
            {
              firstNames: "Rico",
              lastName: "Welch",
              status: "inactive",
              appointmentType: "replacement",
              signedAt: "2022-12-19T09:12:59Z",
            },
            {
              firstNames: "Anne",
              lastName: "Rice",
              status: "active",
              appointmentType: "replacement",
              signedAt: "2022-12-19T07:18:59Z",
              cannotMakeJointDecisions: true,
            },
          ],
          trustCorporations: [
            {
              Name: "Trust Me Ltd.",
              CompanyNumber: "123456789",
              status: "active",
              appointmentType: "original",
            },
            {
              Name: "Trust Me Again Ltd.",
              CompanyNumber: "987654321",
              status: "inactive",
              appointmentType: "replacement",
            },
          ],
          certificateProvider: {
            uid: "e4d5e24e-2a8d-434e-b815-9898620acc71",
            firstNames: "Timothy",
            lastName: "Turner",
            signedAt: "2022-12-18T11:46:24Z",
          },
          restrictionsAndConditions: "Do not do this",
          authorisedSignatory: {
            firstNames: "John",
            lastName: "Signatory",
            signedAt: "2022-12-15T10:30:00Z",
          },
          witnessedByCertificateProviderAt: "2022-12-15T11:00:00Z",
          witnessedByIndependentWitnessAt: "2022-12-15T11:30:00Z",
          independentWitness: {
            firstNames: "Jane",
            lastName: "Witness",
            address: {
              line1: "123 Witness Street",
              line2: "",
              line3: "",
              town: "London",
              postcode: "SW1A 1AA",
              country: "GB",
            },
            email: "jane.witness@example.com",
          },
        },
      }),
      digitalLpas.get("M-DIGI-LPA3-3334", { "opg.poas.lpastore": null }),
      digitalLpas.progressIndicators.feesInProgress("M-DIGI-LPA3-3333"),
      digitalLpas.progressIndicators.feesInProgress("M-DIGI-LPA3-3334"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3334"),
      cases.warnings.empty("333"),
      cases.warnings.empty("1111"),
      cases.tasks.empty("1111"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3333"),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");
  });

  it("shows case information", () => {
    cy.get("h1").contains("Steven Munnell");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3333'] .govuk-tag").contains("Draft");

    cy.contains("PW M-DIGI-LPA3-3334");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3334'] .govuk-tag").contains("Draft");

    cy.contains("PW M-DIGI-LPA3-3335");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3335'] .govuk-tag").contains("Registered");
  });

  it("shows payment information", () => {
    cy.contains("Fees").click();
    cy.contains("£41.00 expected");
  });

  it("shows document information", () => {
    cy.addMock(
      "/lpa-api/v1/lpas/333/documents?type[-][]=Draft&type[-][]=Preview",
      "GET",
      {
        status: 200,
        body: [
          {
            id: 1,
            uuid: "7327f57d-e3d5-4300-95a8-67b3337c7231",
            friendlyDescription: "Mr Test Person - Blank Template",
            direction: "Outgoing",
            createdDate: "24/08/2023 15:27:16",
            systemType: "EP-BB",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "Test",
              surname: "Person",
              personType: "Donor",
            },
          },
          {
            id: 2,
            uuid: "40fa2847-27ae-4976-a93a-9f45ec0a4e98",
            friendlyDescription: "Mr John Doe - Reduced fee evidence",
            direction: "Incoming",
            createdDate: "15/05/2023 11:09:28",
            receivedDateTime: "15/05/2023 11:09:28",
            type: "Application Related",
            subType: "Reduced fee request evidence",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "John",
              surname: "Doe",
              personType: "Correspondent",
            },
          },
        ],
      },
    );
    cy.contains("Documents").click();

    cy.contains("Mr Test Person - Blank Template");
    cy.contains("[OUT]");
    cy.contains("24 August 2023");
    cy.contains("EP-BB");

    cy.contains("Mr John Doe - Reduced fee evidence");
    cy.contains("[IN]");
    cy.contains("15 May 2023");
    cy.contains("Application Related");
    cy.contains("Reduced fee request evidence");
  });

  it("shows task table", () => {
    cy.get(
      "table[data-role=tasks-table] [data-role=tasks-table-header] tr th",
    ).should((elts) => {
      expect(elts).to.contain("Tasks");
      expect(elts).to.contain("Due date");
      expect(elts).to.contain("Actions");
    });
    cy.get(
      "table[data-role=tasks-table] tr[data-role=tasks-table-task-row]",
    ).should((elts) => {
      expect(elts).to.have.length(3);
      expect(elts).to.contain("Review reduced fee eligibility");
      expect(elts).to.contain("Review application correspondence");
      expect(elts).to.contain("Another task");
      expect(elts).to.contain("Reassign task");
    });
  });

  it("shows warnings list", () => {
    cy.addMock("/lpa-api/v1/cases/333/warnings", "GET", {
      status: 200,
      body: [
        {
          id: 44,
          warningType: "Court application in progress",
          warningText: "Court notified",
          dateAdded: "24/08/2022 13:13:13",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
          ],
        },
        {
          id: 22,
          warningType: "Complaint Received",
          warningText: "Complaint from donor",
          dateAdded: "12/12/2023 12:12:12",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
          ],
        },
        {
          id: 24,
          warningType: "Donor Deceased",
          warningText: "Advised of donor death",
          dateAdded: "05/01/2022 10:10:00",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
            { uId: "M-DIGI-LPA3-6666", caseSubtype: "personal-welfare" },
          ],
        },
      ],
    });

    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get(".app-caseworker-summary > div:nth-child(2) li").should((elts) => {
      expect(elts).to.have.length(3);

      expect(elts[0]).to.contain("Donor Deceased");
      expect(elts[0]).to.contain(
        "this case, PA M-DIGI-LPA3-5555 and PW M-DIGI-LPA3-6666",
      );

      expect(elts[1]).to.contain("Complaint Received");
      expect(elts[1]).to.contain("this case and PA M-DIGI-LPA3-5555");

      expect(elts[2]).to.contain("Court application in progress");
      expect(elts[2]).not.to.contain("this case");
    });
  });

  it("creates a warning via case actions", () => {
    cy.contains(".govuk-button", "Case actions").click();
    cy.contains("Create a warning").click();
    cy.url().should("include", "/create-warning?id=33");
    cy.get("#f-warningType").select("Complaint Received");
    cy.get("#f-warningText").type("Be warned!");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Warning created");
    cy.get("h1").contains("Steven Munnell");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("shows lpa details from store when status is Processing", () => {
    cy.contains("Attorneys (4)");
    cy.contains("Replacement attorneys (3)");
    cy.contains("Removed attorneys (1)");
    cy.contains("Notified people (0)");
    cy.contains("Correspondent");

    cy.contains("Review and confirm if severance is required").should(
      "not.exist",
    );
  });

  it("allows changing lpa decisions", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/decisions", "PUT", {
      status: 200,
    });

    cy.contains("Decisions").click();
    cy.get("#f-update-decisions").click();
    cy.contains("Jointly for some").click();
    cy.get("#f-howAttorneysMakeDecisionsDetails").type("This way");
    cy.contains("Continue").click();
    cy.contains("Update saved");
  });

  it("shows channel for donor", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details").then(() => {
      Cypress.$("span:contains('Donor')").closest("button")[0].click();
      cy.contains("Online");
    });
  });

  it("shows due date", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details").then(() => {
      cy.contains("Registration due: 1 December 2023");
    });
  });

  it("shows attorney signed on date and label if set", () => {
    cy.contains("Attorneys (4)")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("Greenwood")
          .parents(".govuk-summary-list")
          .contains("Signed on")
          .should("not.exist");
        cy.contains("McSpolo")
          .parents(".govuk-summary-list")
          .should("contain", "Signed on")
          .and("contain", "20 December 2022");
      });
  });

  it("shows replacement attorney signed on date and label if set", () => {
    cy.contains("Replacement attorneys (3)")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("Guinea")
          .parents(".govuk-summary-list")
          .contains("Signed on")
          .should("not.exist");
        cy.contains("Welch")
          .parents(".govuk-summary-list")
          .should("contain", "Signed on")
          .and("contain", "19 December 2022");
      });
  });

  it("shows certificate provider signed on date, label and change link", () => {
    cy.contains("Certificate provider")
      .click()
      .parents(".govuk-accordion__section")
      .should("contain", "Signed on")
      .and("contain", "18 December 2022")
      .find("#f-change-certificate-provider-details")
      .should("contain", "Update");
  });

  it("shows application details when store is empty", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3334");

    cy.contains("LPA details").click();
    cy.contains("Application format");
    cy.contains("Paper");
  });

  it("review severance messages appears when review restrictions tasks is open", () => {
    cy.addMock(
      "/lpa-api/v1/cases/333/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [
            {
              id: 2,
              name: "Review restrictions and conditions",
              duedate: "10/12/2023",
              status: "OPEN",
              assignee: { displayName: "Super Team" },
            },
          ],
        },
      },
    );

    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");
    cy.contains("Review and confirm if severance is required");
  });

  it("shows restrictions and conditions", () => {
    cy.contains("button", "Restrictions and conditions").click();

    cy.contains("Do not do this");
  });

  it("shows certificate provider witness only (1 witness) - authorised signatory", () => {
    const lpaMocks = Promise.allSettled([
      digitalLpas.get("M-DIGI-LPA3-3335", {
        "opg.poas.lpastore": {
          witnessedByCertificateProviderAt: "2022-12-15T11:00:00Z",
          authorisedSignatory: {
            firstNames: "John",
            lastName: "Signatory",
            signedAt: "2022-12-15T10:30:00Z",
          },
        },
      }),
      digitalLpas.objections.empty("M-DIGI-LPA3-3335"),
    ]);
    cy.wrap(lpaMocks);

    cy.visit("/lpa/M-DIGI-LPA3-3335/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on behalf of the donor by").should("exist");
        cy.contains("John Signatory").should("exist");

        cy.contains(".govuk-details__summary", "View witness details").click();

        cy.get(".govuk-details__text").within(() => {
          cy.contains("Signed by witness 1 (certificate provider)").should(
            "exist",
          );
          cy.contains("Signed by witness 2").should("not.exist");
        });
      });
  });

  it("shows both witnesses when both are present (2 witnesses) - authorised signatory", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on behalf of the donor by").should("exist");
        cy.contains("John Signatory").should("exist");

        cy.contains(".govuk-details__summary", "View witness details").click();

        cy.get(".govuk-details__text").within(() => {
          cy.contains("Signed by witness 1 (certificate provider)").should(
            "exist",
          );
          cy.contains("Signed by witness 2").should("exist");
          cy.contains("Jane Witness").should("exist");
          cy.contains("123 Witness Street").should("exist");
        });
      });
  });

  it("shows donor signed directly (no signed on behalf)", () => {
    const lpaMocks = Promise.allSettled([
      digitalLpas.get("M-DIGI-LPA3-3335"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3335"),
    ]);
    cy.wrap(lpaMocks);

    cy.visit("/lpa/M-DIGI-LPA3-3335/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on");
        cy.contains("19 December 2022");

        cy.contains("LPA signed on behalf of the donor by").should("not.exist");
        cy.contains("View witness details").should("not.exist");
      });
  });
});
