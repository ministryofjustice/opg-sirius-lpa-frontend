import * as cases from "../mocks/cases";

describe("Manage restrictions form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-6666-6666-6666", "GET", {
      status: 200,
      body: {
        uId: "M-6666-6666-6666",
        "opg.poas.sirius": {
          id: 666,
          uId: "M-6666-6666-6666",
          status: "in-progress",
          caseSubtype: "personal-welfare",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
            severanceStatus: "REQUIRED",
          },
        },
        "opg.poas.lpastore": {
          donor: {
            uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
            firstNames: "James",
            lastName: "Rubin",
            otherNamesKnownBy: "Somebody",
            dateOfBirth: "1990-02-22",
            contactLanguagePreference: "en",
            email: "jrubin@mail.example",
          },
          attorneys: [
            {
              uid: "active-attorney-1",
              firstNames: "Julie",
              lastName: "Rutherford",
              address: {
                line1: "15 Cameron Approach",
                line2: "Nether Collier",
                town: "Worcestershire",
                postcode: "BL2 6DI",
                country: "GB",
              },
              status: "active",
              appointmentType: "original",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              mobile: "0123456789",
              email: "j@example.com",
            },
          ],
          status: "in-progress",
          signedAt: "2024-10-18T11:46:24Z",
          lpaType: "pw",
          channel: "online",
          registrationDate: "2024-11-11",
          restrictionsAndConditions: "I only want ...",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/666", "GET", {
      status: 200,
      body: {
        id: 666,
        uId: "M-6666-6666-6666",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 6,
        },
      },
    });

    cases.warnings.empty("666");

    cy.addMock(
      "/lpa-api/v1/cases/666/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [
            {
              id: 6,
              name: "Review restrictions and conditions",
              duedate: "10/12/2025",
              status: "OPEN",
              assignee: { displayName: "Super Team" },
            },
          ],
        },
      },
    );

    cy.visit("/lpa/M-6666-6666-6666/manage-restrictions");
  });

  it("can be visited from the LPA details manage restrictions link", () => {
    cy.visit("/lpa/M-6666-6666-6666/lpa-details").then(() => {
      cy.get(".govuk-accordion__section-button")
        .contains("Restrictions and conditions")
        .click();
      cy.get("#f-manage-restrictions-conditions").click();
      cy.contains("Manage restrictions and conditions");
      cy.url().should("contain", "/lpa/M-6666-6666-6666/manage-restrictions");
      cy.contains("Manage restrictions and conditions");
      cy.contains("Select an option");
      cy.contains("Donor has provided consent to a severance application");
      cy.contains("Donor has refused severance of restriction and conditions");
    });
  });

  it("can go back to changing the severance required option", () => {
    cy.contains("Change").click();
    cy.url().should("contain", "/lpa/M-6666-6666-6666/manage-restrictions");
    cy.url().should("contain", "action=change-severance-required");
    cy.contains("Manage restrictions and conditions");
    cy.contains("Select an option");
    cy.contains("Severance application is not required");
    cy.contains("Severance application is required");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back").click();
    cy.url().should("contain", "/lpa/M-6666-6666-6666/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-6666-6666-6666/lpa-details");
  });

  it("errors when submitting without selecting an option", () => {
    cy.contains("Save and exit").click();
    cy.contains("Please select an option");
  });

  it("redirects when severance application is not required", () => {
    cy.addMock("/lpa-api/v1/tasks/6/mark-as-completed", "PUT", {
      status: 200,
    });
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-6666-6666-6666/severance-status",
      "PUT",
      {
        status: 204,
      },
    );
    cy.contains("Change").click();
    cy.contains("Severance application is not required").click();
    cy.contains("Confirm").click();
    cy.url().should("contain", "/lpa/M-6666-6666-6666/lpa-details");
  });

  it("redirects when severance application is required", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-6666-6666-6666/severance-status",
      "PUT",
      {
        status: 204,
      },
    );
    cy.contains("Change").click();
    cy.contains("Severance application is required").click();
    cy.contains("Confirm").click();
    cy.url().should("contain", "/lpa/M-6666-6666-6666/lpa-details");
  });

  it("Ongoing severance application message appears when severance status is required", () => {
    cy.visit("/lpa/M-6666-6666-6666/lpa-details");
    cy.contains("Ongoing severance application");
  });

  it("Previous not required selection shown in form", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-6666-6666-6668", "GET", {
      status: 200,
      body: {
        uId: "M-6666-6666-6668",
        "opg.poas.sirius": {
          id: 888,
          uId: "M-6666-6666-6668",
          status: "in-progress",
          caseSubtype: "personal-welfare",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: {
            id: 88,
          },
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
            severanceStatus: "NOT_REQUIRED",
          },
        },
      },
    });

    cases.warnings.empty("888");

    cy.addMock("/lpa-api/v1/cases/888", "GET", {
      status: 200,
      body: {
        id: 888,
        uId: "M-6666-6666-6668",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 6,
        },
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/888/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );
    cy.visit("/lpa/M-6666-6666-6668/manage-restrictions").then(() => {
      cy.contains("Severance application required:");
      cy.contains("No");

      cy.contains("Severance application is required");
    });
  });

  it("can be visited from the LPA details manage restrictions link when donor has given consent", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-6666-6666-6669", "GET", {
      status: 200,
      body: {
        uId: "M-6666-6666-6669",
        "opg.poas.sirius": {
          id: 888,
          uId: "M-6666-6666-6669",
          status: "in-progress",
          caseSubtype: "personal-welfare",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: {
            id: 88,
          },
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
            severanceStatus: "REQUIRED",
            severanceApplication: {
              hasDonorConsented: true,
            },
          },
        },
      },
    });

    cases.warnings.empty("888");

    cy.addMock("/lpa-api/v1/cases/888", "GET", {
      status: 200,
      body: {
        id: 888,
        uId: "M-6666-6666-6669",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 6,
        },
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/888/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );
    cy.visit("/lpa/M-6666-6666-6669/manage-restrictions").then(() => {
      cy.contains("Manage restrictions and conditions");
      cy.url().should("contain", "/lpa/M-6666-6666-6669/manage-restrictions");
      cy.contains("Manage restrictions and conditions");
      cy.contains("Date court order made");
      cy.contains("Date court order issued");
      cy.contains(
        "Has severance of the restrictions and conditions been ordered?",
      );
    });
  });
});
